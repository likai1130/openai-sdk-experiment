package assistants

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Threads interface {
	Thread
	ThreadRun
}
type Thread interface {
	CreateThread(ctx context.Context, content string) (openai.Thread, error)
	GetThread(ctx context.Context, threadId string) (openai.Thread, error)
	DeleteThread(ctx context.Context, threadId string) (openai.ThreadDeleteResponse, error)
	ModifyThread(ctx context.Context, threadId string) (openai.Thread, error)
}
type thread struct {
	cli *openai.Client
}

// CreateThread TODO Metadata,FileIDs
func (t *thread) CreateThread(ctx context.Context, content string) (openai.Thread, error) {
	return t.cli.CreateThread(ctx, openai.ThreadRequest{
		Messages: []openai.ThreadMessage{
			{
				Role:    "user",
				Content: content,
				//	FileIDs:  nil,
				//	Metadata: nil,
			},
		},
		//Metadata: nil,
	})
}

func (t *thread) GetThread(ctx context.Context, threadId string) (openai.Thread, error) {
	return t.cli.RetrieveThread(ctx, threadId)
}

func (t *thread) DeleteThread(ctx context.Context, threadId string) (openai.ThreadDeleteResponse, error) {
	return t.cli.DeleteThread(ctx, threadId)
}

// ModifyThread todo Metadata
func (t *thread) ModifyThread(ctx context.Context, threadId string) (openai.Thread, error) {
	return t.cli.ModifyThread(ctx, threadId, openai.ModifyThreadRequest{
		Metadata: nil,
	})
}

func NewThread(client *openai.Client) Threads {
	return &thread{
		cli: client,
	}
}
