package agent

import (
	"context"
	"fmt"

	"github.com/renantatsuo/james-bond/internal/agent/tools"
)

type Client interface {
	SendMessage(ctx context.Context, input []string, model string) (string, error)
	SetTools(tools []tools.Tool)
}

type Agent struct {
	client  Client
	history []string
	tools   []tools.Tool
}

func New(client Client) *Agent {
	return &Agent{
		client:  client,
		history: []string{},
		tools: []tools.Tool{
			tools.ReadFile,
			tools.MyName,
			tools.ListFiles,
		},
	}
}

func (a *Agent) SendUserMessage(ctx context.Context, message, model string) (string, error) {
	a.client.SetTools(a.tools)
	a.history = append(a.history, message)
	response, err := a.client.SendMessage(ctx, a.history, model)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}
	// is this correct? i guess it should have a message type: LLM | User?
	a.history = append(a.history, response)

	return response, nil
}
