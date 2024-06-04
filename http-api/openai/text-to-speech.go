package openai

import (
	"context"
	"fmt"
	"io"

	"github.com/charmbracelet/log"
	goOpenAi "github.com/sashabaranov/go-openai"
)

func TextToSpeech(text string) (response []byte, err error) {
	log.Info("Text to speech called")

	client, err := Client()
	if err != nil {
		log.Error("Error creating client", "err", err)
		return response, err
	}

	req := goOpenAi.CreateSpeechRequest{
		Model:          goOpenAi.TTSModel1,
		Input:          text,
		Voice:          goOpenAi.VoiceNova,
		ResponseFormat: goOpenAi.SpeechResponseFormatWav,
	}

	res, err := client.CreateSpeech(context.Background(), req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return response, err
	}

	buffer, err := io.ReadAll(res)
	if err != nil {
		log.Error("Text to speech to buffer error", "err", err)
	}

	return buffer, nil
}
