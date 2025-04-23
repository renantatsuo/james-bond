package agent

import (
	"context"
)

type Client interface {
	SendMessage(ctx context.Context, input []string, model string) (string, error)
}

type Agent struct {
	client  Client
	history []string
}

func New(client Client) *Agent {
	return &Agent{
		client:  client,
		history: []string{},
	}
}

func (a *Agent) SendUserMessage(ctx context.Context, message, model string) (string, error) {
	a.history = append(a.history, message)
	response, err := a.client.SendMessage(ctx, a.history, model)
	if err != nil {
		return "", err
	}
	a.history = append(a.history, response)

	return response, nil
}
