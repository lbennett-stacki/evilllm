package openai

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	openai "github.com/sashabaranov/go-openai"
)

func Client() (response *openai.Client, err error) {
	log.Info("openai Client called")

	token := os.Getenv("OPENAI_API_TOKEN")

	if token == "" {
		err := errors.New("OPENAI_API_TOKEN is not set")
		log.Error("Error creating OpenAI client", "err", err)
		return nil, err
	}

	client := openai.NewClient(token)

	return client, nil
}
