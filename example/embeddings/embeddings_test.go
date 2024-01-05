package embeddings

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	. "openai-sdk-experiment/example/client"
	"testing"
)

func TestEmbeddings(t *testing.T) {
	ctx := context.Background()
	res, err := Client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{
			"Israel Gaza: Three-year-old-twins among hostages released by Hamas", //以色列-加沙：哈马斯释放的人质中有一对三岁的双胞胎
		},
		Model:          openai.AdaEmbeddingV2,
		EncodingFormat: openai.EmbeddingEncodingFormatBase64,
	})
	if err != nil {
		panic(err)
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}
