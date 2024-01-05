package text

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	. "openai-sdk-experiment/example/client"
	"os"
	"testing"
)

const encoding = "cl100k_base" //文生文的编码方式

/*
*
流式测试
*/
func TestAIFile(t *testing.T) {
	file, err := os.Open("./ai_text_test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	text := string(b)
	tokens := getTokenByEncoding(text, encoding)

	if tokens > 4096 {
		tokens = 4096
	}

	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo1106,
		MaxTokens: tokens,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		Stream: true,
	}
	stream, err := Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	fmt.Println("在耗尽速率限制之前允许的最大请求数。", stream.GetRateLimitHeaders().LimitRequests)              //样本值：60
	fmt.Println("在耗尽速率限制之前允许的最大令牌数", stream.GetRateLimitHeaders().LimitTokens)                 //150000
	fmt.Println("在耗尽速率限制之前允许的剩余请求数", stream.GetRateLimitHeaders().RemainingRequests)           //59
	fmt.Println("在耗尽速率限制之前允许的剩余令牌数量。", stream.GetRateLimitHeaders().RemainingTokens)           //149984
	fmt.Println("速率限制（基于请求）重置为其初始状态的时间。", stream.GetRateLimitHeaders().ResetRequests.String()) //1s
	fmt.Println("速率限制（基于令牌）重置为其初始状态的时间。", stream.GetRateLimitHeaders().ResetTokens.String())   //6m0s
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}

// getTokenByEncoding
func getTokenByEncoding(text string, encoding string) (num_tokens int) {
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf(": %v", err)
		return
	}
	token := tke.Encode(text, nil, nil)
	return len(token)
}
