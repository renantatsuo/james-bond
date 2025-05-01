package agent

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/renantatsuo/james-bond/internal/agent/tools"
	"github.com/swaggest/jsonschema-go"
)

type OpenAIClient struct {
	client  openai.Client
	tools   []openai.ChatCompletionToolParam
	toolMap map[string]tools.ToolFn
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIClient{
		client:  client,
		tools:   []openai.ChatCompletionToolParam{},
		toolMap: map[string]tools.ToolFn{},
	}
}

func (a *OpenAIClient) SendMessage(ctx context.Context, input []string, model string) (string, error) {
	messages := mapMessages(input)

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    getModel(model),
		Tools:    a.tools,
	}

	completion, err := a.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		return completion.Choices[0].Message.Content, nil
	}

	for len(toolCalls) > 0 {
		params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
		for _, toolCall := range toolCalls {
			toolRes, err := a.handleTool(toolCall)
			if err != nil {
				return "", fmt.Errorf("failed to handle tool call %q: %v", toolCall.Function.Name, err)
			}

			params.Messages = append(params.Messages, openai.ToolMessage(toolRes, toolCall.ID))
		}

		completion, err = a.client.Chat.Completions.New(ctx, params)
		if err != nil {
			return "", fmt.Errorf("failed to send tool message: %v", err)
		}
		toolCalls = completion.Choices[0].Message.ToolCalls
	}

	return completion.Choices[0].Message.Content, nil
}

func (a *OpenAIClient) SetTools(tools []tools.Tool) {
	if len(a.tools) > 0 {
		return
	}
	for _, tool := range tools {
		openaiTool := openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.Opt(tool.Description),
				Parameters:  mapToolArgs(tool.Args),
			},
		}
		a.tools = append(a.tools, openaiTool)
		a.toolMap[tool.Name] = tool.Fn
	}
}

func (a *OpenAIClient) handleTool(toolCall openai.ChatCompletionMessageToolCall) (string, error) {
	fn, found := a.toolMap[toolCall.Function.Name]
	if !found {
		return "", fmt.Errorf("tool function %q not found", toolCall.Function.Name)
	}
	args := []byte(toolCall.Function.Arguments)
	return fn(args)
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

func mapMessages(messages []string) []openai.ChatCompletionMessageParamUnion {
	var chatMessages []openai.ChatCompletionMessageParamUnion
	for _, msg := range messages {
		chatMessages = append(chatMessages, openai.UserMessage(msg))
	}
	return chatMessages
}

func mapToolArgs(schema map[string]jsonschema.SchemaOrBool) openai.FunctionParameters {
	if len(schema) == 0 {
		return nil
	}

	params := openai.FunctionParameters{
		"type":       "object",
		"properties": schema,
	}

	return params
}
