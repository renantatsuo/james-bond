package agent

import (
	"context"
	"log"
	"slices"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIClient{client: client}
}

func (a *OpenAIClient) SendMessage(ctx context.Context, input []string, model string) (string, error) {
	messages := mapMessages(input)
	completion, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    getModel(model),
	})

	if err != nil {
		return "", err
	}
	message := transformModel(*completion)
	return message, nil
}

func getModel(model string) openai.ChatModel {
	supportedModels := []openai.ChatModel{
		"gpt-4.1-nano",
	}
	if slices.Contains(supportedModels, model) {
		return model
	} else {
		log.Printf("Unsupported model %q. Using %q as fallback", model, supportedModels[0])
		return supportedModels[0]
	}
}

func transformModel(completion openai.ChatCompletion) string {
	return completion.Choices[0].Message.JSON.Content.Raw()
}

func mapMessages(messages []string) []openai.ChatCompletionMessageParamUnion {
	var chatMessages []openai.ChatCompletionMessageParamUnion
	for _, msg := range messages {
		chatMessages = append(chatMessages, openai.UserMessage(msg))
	}
	return chatMessages
}
