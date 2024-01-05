package assistants

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Assistants interface {
	CreateAssistant(ctx context.Context, assistant AssistantRequest) (openai.Assistant, error)
	GetAssistant(ctx context.Context, assistantId string) (openai.Assistant, error)
	ListAssistants(ctx context.Context, pageTools PageTools) (openai.AssistantsList, error)
	DeleteAssistant(ctx context.Context, assistantId string) (openai.AssistantDeleteResponse, error)
	ModifyAssistant(ctx context.Context, assistantId string, assistant AssistantRequest) (openai.Assistant, error)
}

type AssistantRequest struct {
	Name         string //助手名称
	Description  string // 描述
	Instructions string // 详细说明
	//	Tools        []openai.AssistantTool
}

type PageTools struct {
	Limit  int
	Order  string
	After  string
	Before string
}

type assistants struct {
	cli *openai.Client
}

func (a *assistants) CreateAssistant(ctx context.Context, assistant AssistantRequest) (openai.Assistant, error) {
	return a.cli.CreateAssistant(ctx, openai.AssistantRequest{
		Model:        openai.GPT3Dot5Turbo1106,
		Name:         &assistant.Name,
		Description:  &assistant.Description,
		Instructions: &assistant.Instructions,
		Tools: []openai.AssistantTool{
			{
				Type: openai.AssistantToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "get_current_weather",
					Description: "Get the current weather",
					Parameters:  map[string]any{},
				},
			},
		},
	})
}

func (a *assistants) GetAssistant(ctx context.Context, assistantId string) (openai.Assistant, error) {
	return a.cli.RetrieveAssistant(ctx, assistantId)
}

func (a *assistants) ListAssistants(ctx context.Context, pageTools PageTools) (openai.AssistantsList, error) {
	return a.cli.ListAssistants(ctx, &pageTools.Limit, &pageTools.Order, &pageTools.After, &pageTools.Before)
}

func (a *assistants) DeleteAssistant(ctx context.Context, assistantId string) (openai.AssistantDeleteResponse, error) {
	return a.cli.DeleteAssistant(ctx, assistantId)
}

func (a *assistants) ModifyAssistant(ctx context.Context, assistantId string, assistant AssistantRequest) (openai.Assistant, error) {
	return a.cli.ModifyAssistant(ctx, assistantId, openai.AssistantRequest{
		Model:        openai.GPT3Dot5Turbo1106,
		Name:         &assistant.Name,
		Description:  &assistant.Description,
		Instructions: &assistant.Instructions,
		//Tools:        assistant.Tools,
	})
}

func NewAssistants(client *openai.Client) Assistants {
	return &assistants{
		cli: client,
	}
}
