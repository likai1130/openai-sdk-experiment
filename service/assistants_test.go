package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"io"
	"log"
	"net/http"
	"openai-sdk-experiment/service/assistants"
	"os"
	"testing"
	"time"
)

var assistantId = "asst_HUMjqIbCJ2DXZsqnuPaFIhwY" //数学导师助手
var threadId = "thread_Wk60Uo1CaJe5IoN47qy0d4jm"  //数学导师线程
var cli = createOpenAIClient(os.Getenv("gptkey"))

func TestAssistants(t *testing.T) {
	ctx := context.Background()
	// 用线程发送信息
	thread := assistants.NewThread(cli)
	/*run, err := thread.CreateAndRun(ctx, assistantId, "再帮我解2x+11=21这个方程式")
	if err != nil {
		log.Println("创建信息出错 err=", err)
		return
	}
	getThread, err := thread.GetThread(ctx, run.ThreadID)
	if err != nil {
		log.Println("查看线程出错 err=", err)
		return
	}
	bytes, err := json.Marshal(getThread)
	if err != nil {
		log.Println("getThread json err", err)
		return
	}
	log.Printf("thread 详情 = %v \n", string(bytes))*/
	newMessage := assistants.NewMessage(cli)
	_, err := newMessage.CreateMessage(ctx, threadId, "计算4x + 5 = 20这个方程")
	if err != nil {
		log.Println("createMessage  err", err)
		return
	}
	run, err := thread.Run(ctx, threadId, assistantId)
	if err != nil {
		log.Println("Run  err", err)
		return
	}

	for {
		getRun, err := thread.GetRun(ctx, run.ThreadID, run.ID)
		if err != nil {
			log.Println("for getRun线程出错 err=", err)
			return
		}
		log.Println("当前执行状态为", getRun.Status)
		if getRun.Status == openai.RunStatusCompleted {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 查信息列表
	message := assistants.NewMessage(cli)
	messages, err := message.ListMessages(ctx, run.ThreadID, assistants.PageTools{
		Limit: 100,
		Order: "desc",
	})
	if err != nil {
		log.Println("查询list message错误了，err = ", err)
		return
	}
	marshal, err := json.Marshal(messages)
	if err != nil {
		log.Println("json err = ", err)
		return
	}
	log.Println("消息列表：", string(marshal))

}

func TestName(t *testing.T) {
	ctx := context.Background()
	thread, err := assistants.NewThread(cli).GetThread(ctx, "thread_Wk60Uo1CaJe5IoN47qy0d4jm")
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(thread)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

/*
TestTools 测试助手工具

asst_HUMjqIbCJ2DXZsqnuPaFIhwY
thread_Wk60Uo1CaJe5IoN47qy0d4jm
更新助手，上传代码文件
*/
func createWeatherAss(ctx context.Context) openai.Assistant {
	name := "天气预报员"
	instructions := "你是一个天气机器人。使用提供的功能来回答问题。"
	description := "你是一个天气机器人。使用提供的功能来回答问题。"

	// 1. 创建助手
	response, err := client.CreateAssistant(ctx, openai.AssistantRequest{
		Model:        openai.GPT3Dot5Turbo1106,
		Name:         &name,
		Description:  &description,
		Instructions: &instructions,
		Tools: []openai.AssistantTool{
			{
				Type: openai.AssistantToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name: "get_current_weather",
					Parameters: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"location": {
								Type:        jsonschema.String,
								Description: "The city and state, e.g. San Francisco, CA",
							},
							"unit": {
								Type: jsonschema.String,
								Enum: []string{"celsius", "fahrenheit"},
							},
						},
						Required: []string{"location"},
					},
				},
			},
		},
		//FileIDs: nil,
	})

	if err != nil {
		panic(err)
	}
	return response
}

func createWeatherThread(ctx context.Context, assId string) openai.Run {
	mdoel := "gpt-3.5-turbo-1106"
	response, err := cli.CreateThreadAndRun(ctx, openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{
			AssistantID: assId,
			Model:       &mdoel,
		},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    "user",
					Content: "请告诉我河南省郑州市当前的天气情况",
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(response.ThreadID)
	return response
}

func listMessage(ctx context.Context, threadId string) openai.MessagesList {
	messages, err := cli.ListMessage(ctx, threadId, nil, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	return messages
}

/*
*
asst_wpUPjUfDhb9LK9VRDdRT9Y2h
run_QYfzAf53RJB6vaakfr0IO0fr
*/
func TestAss(t *testing.T) {
	thId := "thread_bC7dn4cnBsFUJUepkmwSherd"
	//ruId := "run_QYfzAf53RJB6vaakfr0IO0fr"
	assId := "asst_wpUPjUfDhb9LK9VRDdRT9Y2h"
	ctx := context.Background()
	run := NewAssistants().AddAndRun(ctx, RunRequest{
		AssistantsId: assId,
		ThreadId:     thId,
		Content:      "请告诉我郑州市的当前天气情况",
	})
	if len(run.Msg) != 0 {
		panic(run.Msg)
	}
	/*	bytes, _ := json.Marshal(run.Data)
		r := openai.Run{}
		_ = json.Unmarshal(bytes, &r)
		fmt.Println("线程", r.ThreadID)
	*/
	message := listMessage(ctx, thId)
	marshal, _ := json.Marshal(message)
	fmt.Println(string(marshal))
}

/*
*
期刊助手id= asst_0znUJEBvKz9jLt34qwEcb4QN
期刊线程id= thread_9Nq2DanJBAHqNJjGDp93ngLL
*/
func TestKuggaAss(t *testing.T) {
	assId := "asst_0znUJEBvKz9jLt34qwEcb4QN"
	thId := "thread_9Nq2DanJBAHqNJjGDp93ngLL"
	ctx := context.Background()
	/*ass := createKuggaAss(ctx)
	fmt.Println("期刊助手id= ", ass.ID)*/
	run := NewAssistants().AddAndRun(ctx, RunRequest{
		AssistantsId: assId,
		ThreadId:     thId,
		Content:      "请帮我查询id为421的文章详情",
	})
	if len(run.Msg) != 0 {
		panic(run.Msg)
	}
	/*bytes, _ := json.Marshal(run.Data)
	r := openai.Run{}
	_ = json.Unmarshal(bytes, &r)
	fmt.Println("线程", r.ThreadID)
	*/
	message := listMessage(ctx, thId)
	marshal, _ := json.Marshal(message)
	fmt.Println(string(marshal))
}

func createKuggaAss(ctx context.Context) openai.Assistant {
	name := "期刊管理员"
	instructions := "你是一个期刊管理员。根据我提供的id，为我获取对应的文章内容，文章内容的字段在是content字段。"
	description := "你是一个期刊管理员。根据我提供的id，为我获取对应的文章内容，文章内容的字段在是content字段。"
	// 1. 创建助手
	response, err := client.CreateAssistant(ctx, openai.AssistantRequest{
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

	if err != nil {
		panic(err)
	}
	return response
}

/*
*
写一个函数
*/
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
	fmt.Println(response)
	return response
}
