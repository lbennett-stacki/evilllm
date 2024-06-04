package routes

import (
	"encoding/json"
	"evilllm-http-api/openai"
	"evilllm-http-api/upload"
	"fmt"
	"net/http"
	"os"

	goOpenAi "github.com/sashabaranov/go-openai"

	"github.com/charmbracelet/log"
)

const (
	shouldMockUpload = false
	shouldMockTTS    = false
	shouldMockSTT    = false
	shouldMockChat   = false
)

var GAME_MESSAGES []goOpenAi.ChatCompletionMessage

func CommunicateHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Communicate handler called")

	var fileName string

	if shouldMockUpload {
		fileName = "audio.wav"
	} else {
		file, err := upload.UploadFile(r)
		if err != nil {
			log.Error("Speech to text error", "err", err)
			fmt.Fprint(w, "error")
			return
		}
		fileName = file
	}

	var text string
	if shouldMockSTT {
		text = "You have enabled STT mocking in the Evil LLM API server."
	} else {
		log.Info("Calling speech to text")
		humanText, err := openai.SpeechToText(upload.UploadsPath(fileName))
		if err != nil {
			log.Error("Speech to text error", "err", err)
			fmt.Fprint(w, "error")
			return
		}
		text = humanText
	}

	var replyText string

	if shouldMockChat {
		replyText = "You have enabled chat mocking in the Evil LLM API server."
	} else {
		log.Info("Calling chat for reply", "text", text)
		newMessages, err := openai.Chat(text, GAME_MESSAGES)
		if err != nil {
			log.Error("Error getting chat completion reply text", "err", err)
			fmt.Fprint(w, "error")
			return
		}
		GAME_MESSAGES = newMessages

		var reply openai.ReplyJsonSchema

		replyContent := newMessages[len(newMessages)-1].Content

		json.Unmarshal([]byte(replyContent), &reply)

		log.Info("Unmarhalled a reply", "reply", reply, "from replyContent", replyContent)

		if reply.ReplyText == "" {
			log.Error("Error handling evil robot response, no ReplyText found in response", "reply", reply, "replyContent", replyContent)
			fmt.Fprint(w, "error")
		}

		if reply.IsChallengeComplete {
			log.Info("User has won!!!!!")
		}

		replyText = reply.ReplyText
	}

	var reply []byte

	if shouldMockTTS {
		replyBuffer, err := os.ReadFile(upload.GeneratedPath(fileName))
		if err != nil {
			log.Error("Read tts mock error", "err", err)
			fmt.Fprint(w, "error")
			return
		}
		reply = replyBuffer
	} else {
		log.Info("Calling tts for reply", "replyText", replyText)
		replyBuffer, err := openai.TextToSpeech(replyText)
		if err != nil {
			log.Error("Text to speech error", "err", err)
			fmt.Fprint(w, "error")
			return
		}
		reply = replyBuffer
	}

	log.Debug("Saving reply mp3", "filename", fileName)
	err := os.WriteFile(upload.GeneratedPath(fileName), reply, 0644)
	if err != nil {
		log.Error("Text to speech to file error", "err", err)
	}

	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(reply)))
	w.Header().Set("Content-Disposition", "inline; filename=\"output.wav\"")

	log.Debug("Writing reply byte array", "filename", fileName)
	_, err = w.Write(reply)
	if err != nil {
		log.Debug("Error writing response:", "err", err)
	}
}
