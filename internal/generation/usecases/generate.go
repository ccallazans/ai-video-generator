package usecases

import (
	"log"
	"os"
	"os/exec"

	"github.com/ccallazans/ai-video-generator/internal/generation/processes"
)

func Generate(message string) (string, error) {

	// initPythonEnvironment() // For debug

	tempDir, err := os.MkdirTemp("", "ai-video-generator")
	if err != nil {
		log.Println("Failed to create temporary directory: ", err.Error())
		return "", err
	}
	defer os.RemoveAll(tempDir)

	textProcess := processes.NewLocalTextGeneration()
	generatedText, err := textProcess.Execute(message)
	if err != nil {
		return "", err
	}

	speechProcess := processes.NewLocalSpeechGeneration(tempDir)
	speechFilename, err := speechProcess.Execute(generatedText)
	if err != nil {
		return "", err
	}

	videoProcess := processes.NewLocalVideoGeneration(tempDir, speechFilename)
	finalVideo, err := videoProcess.Execute("")
	if err != nil {
		return "", err
	}

	return finalVideo, nil
}

func initPythonEnvironment() error {
	cmd := exec.Command("bash", "-c", "source venv/bin/activate")
	otp, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error initializing python environment")
		log.Println(string(otp))
		return err
	}

	return nil
}
