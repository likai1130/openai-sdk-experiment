package service

import (
	"github.com/sashabaranov/go-openai"
	"os"
	"sync"
)

var client *openai.Client
var once sync.Once

func init() {
	createOpenAIClient(os.Getenv("OPENAI_KEY"))
}
func createOpenAIClient(token string) *openai.Client {
	once.Do(func() {
		client = openai.NewClient(token)
	})
	return client
}
