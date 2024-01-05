package assistants

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Messages interface {
	CreateMessage(ctx context.Context, threadId, msg string) (openai.Message, error)
	GetMessage(ctx context.Context, threadId, messageId string) (openai.Message, error)
	ListMessages(ctx context.Context, threadId string, pageTools PageTools) (openai.MessagesList, error)
}

type message struct {
	cli *openai.Client
}

func (m *message) CreateMessage(ctx context.Context, threadId, msg string) (openai.Message, error) {
	return m.cli.CreateMessage(ctx, threadId, openai.MessageRequest{
		Role:    "user",
		Content: msg,
	})
}

func (m *message) GetMessage(ctx context.Context, threadId, messageId string) (openai.Message, error) {
	return m.cli.RetrieveMessage(ctx, threadId, messageId)
}

func (m *message) ListMessages(ctx context.Context, threadId string, pageTools PageTools) (openai.MessagesList, error) {
	return m.cli.ListMessage(ctx, threadId, &pageTools.Limit, &pageTools.Order, &pageTools.After, &pageTools.Before)
}

func NewMessage(client *openai.Client) Messages {
	return &message{cli: client}
}
