package openai

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	openai "github.com/sashabaranov/go-openai"
)

func SpeechToText(filePath string) (response string, err error) {
	log.Info("Speech to text called")

	client, err := Client()
	if err != nil {
		log.Error("Error creating client", "err", err)
		return response, err
	}

	log.Info("Speech to text created client")

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filePath,
	}

	log.Info("Speech to text creating transcript")

	resp, err := client.CreateTranscription(context.Background(), req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return response, err
	}

	log.Info("Speech to text returning")

	return resp.Text, nil
}
