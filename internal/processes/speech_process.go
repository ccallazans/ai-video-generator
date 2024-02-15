package processes

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/ccallazans/ai-video-generator/internal/utils"
)

type LocalSpeechGeneration struct {
	tempFolder string
}

func NewLocalSpeechGeneration(tempFolder string) SpeechProcess {
	return &LocalSpeechGeneration{tempFolder: tempFolder}
}

func (p *LocalSpeechGeneration) Execute(command string) (string, error) {
	log.Println("Starting speech process")

	result, err := p.generateSpeech(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *LocalSpeechGeneration) generateSpeech(message string) (string, error) {
	filename := fmt.Sprintf("%s/%s.mp3", p.tempFolder, utils.RandomString())

	args := []string{
		"./pkg/tiktokvoice.py",
		message,
		filename,
	}

	cmd := exec.Command("python", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing script local speech generation: ", err)
		return "", err
	}

	return filename, nil
}
