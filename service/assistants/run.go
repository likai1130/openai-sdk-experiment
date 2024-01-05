package assistants

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type ThreadRun interface {
	CreateAndRun(ctx context.Context, assistantId string, content string) (openai.Run, error)
	Run(ctx context.Context, threadId, assistantId string) (openai.Run, error)
	GetRun(ctx context.Context, threadId, runId string) (openai.Run, error)
	CancelRun(ctx context.Context, threadId, runId string) (openai.Run, error)
	ListRuns(ctx context.Context, threadId string, pageTools PageTools) (openai.RunList, error)
	GetRunStep(ctx context.Context, threadId, runId, stepId string) (openai.RunStep, error)
	ListRunStep(ctx context.Context, threadId, runId string, pageTools PageTools) (openai.RunStepList, error)
	SubmitToolOutputs(ctx context.Context, threadId, runId string, outPuts *openai.SubmitToolOutputs) (openai.Run, error)
}

// CreateAndRun TODO Metadata
func (t *thread) CreateAndRun(ctx context.Context, assistantId string, content string) (openai.Run, error) {
	return t.cli.CreateThreadAndRun(ctx, openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{
			AssistantID: assistantId,
		},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    openai.ThreadMessageRoleUser,
					Content: content,
					//	FileIDs:  nil,
					//	Metadata: nil,
				},
			},
			//Metadata: nil,
		},
	})
}

// Run TODO Metadata
func (t *thread) Run(ctx context.Context, threadId, assistantId string) (openai.Run, error) {
	return t.cli.CreateRun(ctx, threadId, openai.RunRequest{
		AssistantID: assistantId,
	})
}

func (t *thread) GetRun(ctx context.Context, threadId, runId string) (openai.Run, error) {
	return t.cli.RetrieveRun(ctx, threadId, runId)
}

func (t *thread) CancelRun(ctx context.Context, threadId, runId string) (openai.Run, error) {
	return t.cli.CancelRun(ctx, threadId, runId)
}

func (t *thread) ListRuns(ctx context.Context, threadId string, pageTools PageTools) (openai.RunList, error) {
	return t.cli.ListRuns(ctx, threadId, openai.Pagination{
		Limit:  &pageTools.Limit,
		Order:  &pageTools.Order,
		After:  &pageTools.After,
		Before: &pageTools.Before,
	})
}

func (t *thread) GetRunStep(ctx context.Context, threadId, runId, stepId string) (openai.RunStep, error) {
	return t.cli.RetrieveRunStep(ctx, threadId, runId, stepId)
}

func (t *thread) SubmitToolOutputs(ctx context.Context, threadId, runId string, outPuts *openai.SubmitToolOutputs) (openai.Run, error) {
	var outP []openai.ToolOutput
	for _, outCall := range outPuts.ToolCalls {
		o := openai.ToolOutput{
			ToolCallID: outCall.ID,
			Output:     true,
		}
		outP = append(outP, o)
	}

	return t.cli.SubmitToolOutputs(ctx, threadId, runId, openai.SubmitToolOutputsRequest{
		ToolOutputs: outP,
	})
}

func (t *thread) ListRunStep(ctx context.Context, threadId, runId string, pageTools PageTools) (openai.RunStepList, error) {
	return t.cli.ListRunSteps(ctx, threadId, runId, openai.Pagination{
		Limit:  &pageTools.Limit,
		Order:  &pageTools.Order,
		After:  &pageTools.After,
		Before: &pageTools.Before,
	})
}
