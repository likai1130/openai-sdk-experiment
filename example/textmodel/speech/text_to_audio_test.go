package speech

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	. "openai-sdk-experiment/example/client"
	"os"
	"testing"
)

/*
*
测试从文章内容转语音，alloy声音
*/
func TestSignalArticle(t *testing.T) {
	filePath := "./kugga_article_421.txt"
	text := readFile(filePath)
	model := openai.VoiceAlloy
	if err := speechToMp3(text, "alloy.mp3", model); err != nil {
		log.Printf("alloy 生成错误: err=%v \n", err)
	}
}

/*
*
测试从文章内容转语音，多种声音
*/
func TestBatchArticle(t *testing.T) {
	filePath := "./kugga_article_421.txt"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	text := string(b)

	models := map[string]openai.SpeechVoice{
		"alloy":   openai.VoiceAlloy,
		"echo":    openai.VoiceEcho,
		"fable":   openai.VoiceFable,
		"onyx":    openai.VoiceOnyx,
		"nova":    openai.VoiceNova,
		"shimmer": openai.VoiceShimmer,
	}

	for k, _ := range models {
		v := models[k]
		if err := speechToMp3(text, k+".mp3", v); err != nil {
			log.Printf("%s 生成错误: err=%v \n", k, err)
		}
	}
}
func speechToMp3(text, filePath string, model openai.SpeechVoice) error {
	ctx := context.Background()
	req := openai.CreateSpeechRequest{
		Model: openai.TTSModel1HD,
		Input: text,
		Voice: model,
		//ResponseFormat: "",
		//默认响应格式为“mp3”，但也可以使用“opus”、“aac”或“flac”等其他格式。
		//Opus：用于互联网流媒体和通信，低延迟。AAC：用于数字音频压缩，
		//YouTube、Android、iOS 首选。FLAC：用于无损音频压缩，受到音频爱好者存档的青睐。
		//Speed: 0,
	}

	response, err := Client.CreateSpeech(ctx, req)
	if err != nil {
		return errors.WithMessage(err, "CreateSpeech")
	}
	defer response.Close()

	buf, err := io.ReadAll(response)
	if err != nil {
		return errors.WithMessage(err, "IOReadAll")
	}

	// save buf to file as mp3
	err = os.WriteFile(filePath, buf, 0644)
	if err != nil {
		return errors.WithMessage(err, "WriteFile")
	}
	return nil
}

func readFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
