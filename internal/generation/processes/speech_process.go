package processes

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/ccallazans/ai-video-generator/internal/utils"
)

type SpeechGenerationProcess struct {
	next    Process
	tempDir string
}

func NewSpeechGenerationProcess(tempDir string) *SpeechGenerationProcess {
	return &SpeechGenerationProcess{tempDir: tempDir}
}

func (p *SpeechGenerationProcess) Execute(request interface{}) (interface{}, error) {
	generatedText, ok := request.(string)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	speechFilename, err := p.generateSpeech(generatedText)
	if err != nil {
		return nil, err
	}

	if p.next != nil {
		return p.next.Execute(speechFilename)
	}

	return speechFilename, nil
}

func (p *SpeechGenerationProcess) SetNext(handler Process) {
	p.next = handler
}

func (p *SpeechGenerationProcess) generateSpeech(message string) (string, error) {
	filename := fmt.Sprintf("%s/%s.mp3", p.tempDir, utils.RandomString())

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
