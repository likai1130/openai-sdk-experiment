package client

import (
	"github.com/sashabaranov/go-openai"
	"os"
	"sync"
)

var Client *openai.Client

var once sync.Once

func init() {
	createOpenAIClient(os.Getenv("OPENAI_KEY"))
}

func createOpenAIClient(token string) *openai.Client {
	once.Do(func() {
		Client = openai.NewClient(token)
	})
	return Client
}
