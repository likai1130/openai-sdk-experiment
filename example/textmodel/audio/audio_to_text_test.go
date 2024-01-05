package audio

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	. "openai-sdk-experiment/example/client"
	"testing"
)

/*
*
语音转文字
*/
func TestAudio(t *testing.T) {
	ctx := context.Background()
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "ai.mp4",
	}
	resp, err := Client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}
	fmt.Println(resp.Text)
}
