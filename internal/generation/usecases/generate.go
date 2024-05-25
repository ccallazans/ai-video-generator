package usecases

import (
	"errors"
	"log"
	"os"

	"github.com/ccallazans/ai-video-generator/internal/generation/processes"
)

func Generate(message string) (string, error) {
	tempDir, err := os.MkdirTemp("", "ai-video-generator")
	if err != nil {
		log.Println("Failed to create temporary directory: ", err.Error())
		return "", err
	}
	defer os.RemoveAll(tempDir)

	textProcess := processes.NewTextGenerationProcess()
	speechProcess := processes.NewSpeechGenerationProcess(tempDir)
	videoProcess := processes.NewVideoGenerationProcess(tempDir)

	textProcess.SetNext(speechProcess)
	speechProcess.SetNext(videoProcess)

	result, err := textProcess.Execute(message)
	if err != nil {
		return "", err
	}

	finalVideo, ok := result.(string)
	if !ok {
		return "", errors.New("failed to generate video")
	}

	return finalVideo, nil
}
