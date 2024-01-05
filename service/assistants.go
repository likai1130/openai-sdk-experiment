package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"openai-sdk-experiment/service/assistants"
	"time"
)

type Assistants struct {
	assistants assistants.Assistants
	message    assistants.Messages
	threads    assistants.Threads
}

type RunRequest struct {
	AssistantsId string `json:"assistantsId" binding:"required"`
	ThreadId     string `json:"threadId"`
	Content      string `json:"content" binding:"required"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// AddAndRun 添加信息并且运行它
func (a *Assistants) AddAndRun(ctx context.Context, req RunRequest) Response {
	var err error
	run := openai.Run{}
	if len(req.ThreadId) == 0 {
		run, err = a.threads.CreateAndRun(ctx, req.AssistantsId, req.Content)
		if err != nil {
			return internalServerError(err)
		}
	} else {
		_, err = a.message.CreateMessage(ctx, req.ThreadId, req.Content)
		if err != nil {
			return internalServerError(err)
		}
		run, err = a.threads.Run(ctx, req.ThreadId, req.AssistantsId)
		if err != nil {
			return internalServerError(err)
		}
	}
	for {
		step, err := a.threads.ListRunStep(ctx, run.ThreadID, run.ID, assistants.PageTools{
			Limit: 100,
			Order: "desc",
		})
		if err != nil {
			return internalServerError(errors.WithMessage(err, "ListRunStep出错"))
		}
		marshal, _ := json.Marshal(step)
		fmt.Println("step:", string(marshal))

		getRun, err := a.threads.GetRun(ctx, run.ThreadID, run.ID)
		if err != nil {
			return internalServerError(errors.WithMessage(err, "for getRun线程出错"))
		}
		log.Println("当前执行状态为", getRun.Status)
		if getRun.Status == openai.RunStatusRequiresAction {
			log.Println("当前执行SubmitToolOutputs")
			toolOutputs := getRun.RequiredAction.SubmitToolOutputs
			getRun, err = a.threads.SubmitToolOutputs(ctx, run.ThreadID, run.ID, toolOutputs)
			if err != nil {
				return internalServerError(errors.WithMessage(err, "SubmitToolOutputs 提交tools数据失败"))
			}
		}
		if getRun.Status == openai.RunStatusCompleted {
			return ok(getRun)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// Clear 删除线程
func (a *Assistants) Clear(ctx context.Context, threadId string) Response {
	thread, err := a.threads.DeleteThread(ctx, threadId)
	if err != nil {
		return internalServerError(err)
	}
	return ok(thread)
}

// CreateAssistant 创建助手
func (a *Assistants) CreateAssistant(ctx context.Context, request assistants.AssistantRequest) Response {
	assistant, err := a.assistants.CreateAssistant(ctx, request)
	if err != nil {
		return internalServerError(err)
	}
	return ok(assistant)
}

// ListMessages 消息列表
func (a *Assistants) ListMessages(ctx context.Context, threadId string, tools assistants.PageTools) Response {
	if tools.Limit > 100 {
		return Response{Code: http.StatusBadRequest, Msg: errors.New("分页数量不能大于100").Error()}
	}
	messages, err := a.message.ListMessages(ctx, threadId, tools)
	if err != nil {
		return Response{Code: http.StatusBadRequest, Msg: err.Error()}
	}

	return Response{Code: http.StatusOK, Data: messages}
}

func internalServerError(err error) Response {
	return Response{Code: http.StatusInternalServerError, Msg: err.Error()}
}

func ok(data interface{}) Response {
	return Response{Code: http.StatusOK, Data: data}
}

func NewAssistants() *Assistants {
	return &Assistants{
		assistants: assistants.NewAssistants(client),
		message:    assistants.NewMessage(client),
		threads:    assistants.NewThread(client),
	}
}
