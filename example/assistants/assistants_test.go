package assistants

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"io"
	"log"
	"net/http"
	. "openai-sdk-experiment/example/client"
	"testing"
	"time"
)

type IdObject struct {
	Id int `json:"id"`
}

func (idObj *IdObject) getId(data string) int {
	b := []byte(data)
	json.Unmarshal(b, idObj)
	return idObj.Id

}

func objToStr(data any) string {
	info := ArticleDetails{}
	marshal, _ := json.Marshal(data)
	err := json.Unmarshal(marshal, &info)
	if err != nil {
		log.Println(err)
	}
	return info.DocumentInfo.Content
}

type kuggaResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func getArticlesById(id int) kuggaResponse {
	url := fmt.Sprintf("https://api.kugga.com/api/v1/articles/%d", id)
	defaultClient := http.DefaultClient
	resp, err := defaultClient.Get(url)
	if err != nil {
		return kuggaResponse{http.StatusInternalServerError, "unknow", nil}
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return kuggaResponse{http.StatusInternalServerError, "io.ReadAll unknow", nil}
	}
	response := kuggaResponse{}
	if err := json.Unmarshal(all, &response); err != nil {
		return kuggaResponse{http.StatusInternalServerError, "json.Unmarshal unknow", nil}
	}
	//fmt.Println(response)
	return response
}

func createKuggaAss(ctx context.Context) (openai.Assistant, error) {
	name := "期刊管理员"
	instructions := "你是一个期刊管理员。根据我提供的id，为我获取对应的文章内容，不要修改函数的返回结果。"
	description := "你是一个期刊管理员。根据我提供的id，为我获取对应的文章内容，不要修改函数的返回结果"
	// 1. 创建助手
	return Client.CreateAssistant(ctx, openai.AssistantRequest{
		Model:        openai.GPT3Dot5Turbo1106,
		Name:         &name,
		Description:  &description,
		Instructions: &instructions,
		Tools: []openai.AssistantTool{
			{
				Type: openai.AssistantToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name: "getArticlesById",
					Parameters: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"id": {
								Type:        jsonschema.Integer,
								Description: "这个是文章的id，根据id可以获取到文章的详情",
							},
						},
						Required: []string{"id"},
					},
				},
			},
		},
		//FileIDs: nil,
	})

}

func addAndRun(ctx context.Context, assistantId string, message string) (openai.Run, error) {
	model := "gpt-3.5-turbo-1106"
	return Client.CreateThreadAndRun(ctx, openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{
			AssistantID: assistantId,
			Model:       &model,
		},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    "user",
					Content: message,
				},
			},
		},
	})

}

func run(ctx context.Context, assistantId, threadId string, message string) (openai.Run, error) {
	_, err := Client.CreateMessage(ctx, threadId, openai.MessageRequest{
		Role:    "user",
		Content: message,
	})
	if err != nil {
		return openai.Run{}, err
	}
	return Client.CreateRun(ctx, threadId, openai.RunRequest{
		AssistantID: assistantId,
	})
}

func createThread(ctx context.Context, message string) (response openai.Thread, err error) {
	return Client.CreateThread(ctx, openai.ThreadRequest{
		Messages: []openai.ThreadMessage{
			{
				Role:    "user",
				Content: message},
		},
	})
}

func onceEexec(ctx context.Context, message string) string {
	ass, err := createKuggaAss(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("第一次创建助手：助手id=", ass.ID)
	marshal, err := json.Marshal(ass)
	if err != nil {
		panic(err)
	}
	log.Printf("ass的结果= %s\n", string(marshal))

	response, err := addAndRun(ctx, ass.ID, message)
	if err != nil {
		panic(err)
	}
	log.Println("第一次创建线程：线程id=", response.ThreadID)
	marshal2, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("run 的结果 = %s\n", string(marshal2))

	// 监控状况
	for {
		response, err = Client.RetrieveRun(ctx, response.ThreadID, response.ID)
		if err != nil {
			panic(errors.WithMessage(err, "for getRun线程出错"))
		}
		log.Println("当前执行状态为", response.Status)
		if response.Status == openai.RunStatusRequiresAction {
			log.Println("当前执行SubmitToolOutputs")
			var outputs []openai.ToolOutput
			calls := response.RequiredAction.SubmitToolOutputs.ToolCalls
			for _, c := range calls {
				name := c.Function.Name
				arguments := c.Function.Arguments
				// 判断并且读取函数
				idObject := IdObject{}
				id := idObject.getId(arguments)
				if name == "getArticlesById" {
					d := getArticlesById(id)

					result := objToStr(d.Data)
					output := openai.ToolOutput{
						ToolCallID: c.ID,
						Output:     result,
					}
					outputs = append(outputs, output)
				}
			}
			if len(outputs) > 0 {
				_, err = subOutPut(ctx, response.ThreadID, response.ID, outputs)
				if err != nil {
					panic(err)
				}

			}
		}
		if response.Status == openai.RunStatusCompleted {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	return response.ThreadID
}

func subOutPut(ctx context.Context,
	threadID string,
	runID string,
	request []openai.ToolOutput) (response openai.Run, err error) {
	return Client.SubmitToolOutputs(ctx, threadID, runID, openai.SubmitToolOutputsRequest{
		ToolOutputs: request,
	})
}

func TestAssistantsArticle(t *testing.T) {
	ctx := context.Background()
	threadId := onceEexec(ctx, "帮我查询id为421的文章详情")
	messages, err := Client.ListMessage(ctx, threadId, nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(messages)
	fmt.Println(string(marshal))
	//run(ctx,assId,threadId,message)
}
