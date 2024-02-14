package usecases

import (
	"log"

	"github.com/ccallazans/ai-video-generator/internal/usecases/processes"
)

func Generate(message string) (string, error) {

	generatedText, err := processes.TextProcess(message)
	if err != nil {
		log.Println(err)
		return "", err
	}

	generatedSpeech, err := processes.SpeechProcess(generatedText)
	if err != nil {
		log.Println(err)
		return "", err
	}

	videos, err := processes.VideoProcess(generatedSpeech)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return videos, nil
}
