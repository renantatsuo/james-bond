package agent

import (
	"context"
	"fmt"

	"github.com/renantatsuo/james-bond/internal/agent/tools"
)

type Client interface {
	SendMessage(ctx context.Context, input []Message, model string) (string, error)
	SetTools(tools []tools.Tool)
}

type Agent struct {
	client  Client
	history []Message
	tools   []tools.Tool
}

type Message struct {
	Type    MessageType
	Content string
}

type MessageType string

const (
	MessageTypeUser MessageType = "User"
	MessageTypeAI   MessageType = "AI"
)

func New(client Client) *Agent {
	return &Agent{
		client:  client,
		history: []Message{},
		tools: []tools.Tool{
			tools.ReadFile,
			tools.MyName,
			tools.ListFiles,
			tools.WriteFile,
		},
	}
}

func (a *Agent) SendUserMessage(ctx context.Context, message Message, model string) (string, error) {
	a.client.SetTools(a.tools)
	a.history = append(a.history, message)
	response, err := a.client.SendMessage(ctx, a.history, model)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	a.history = append(a.history, Message{Type: MessageTypeAI, Content: response})

	return response, nil
}
