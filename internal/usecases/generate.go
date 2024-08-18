package usecases

import (
	"errors"
	"log"
	"os"

	"github.com/ccallazans/ai-video-generator/internal/processes"
)

func Generate(prompt string) (string, error) {
	tempDir, err := os.MkdirTemp("", "ai-video-generator")
	if err != nil {
		log.Println("Failed to create temporary directory: ", err.Error())
		return "", err
	}
	defer os.RemoveAll(tempDir)

	context := &processes.GenerationContext{
		TempDir: tempDir,
		Prompt:  prompt,
	}

	textProcess := processes.NewTextGenerationProcess()
	speechProcess := processes.NewSpeechGenerationProcess()
	videoProcess := processes.NewVideoGenerationProcess()

	textProcess.SetNext(speechProcess)
	speechProcess.SetNext(videoProcess)

	// Execute the pipeline
	result, err := textProcess.Execute(context)
	if err != nil {
		return "", err
	}

	finalVideo, ok := result.(string)
	if !ok {
		return "", errors.New("failed to generate video")
	}

	return finalVideo, nil
}
