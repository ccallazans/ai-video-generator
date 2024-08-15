package processes

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/google/uuid"
)

type SpeechGenerationProcess struct {
	next Process
}

func NewSpeechGenerationProcess() *SpeechGenerationProcess {
	return &SpeechGenerationProcess{}
}

func (p *SpeechGenerationProcess) Execute(request interface{}) (interface{}, error) {
	context, ok := request.(*GenerationContext)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	speechFilename, err := p.generateSpeech(context.Text, context.TempDir)
	if err != nil {
		return nil, err
	}
	context.SpeechFile = speechFilename

	if p.next != nil {
		return p.next.Execute(context)
	}

	return context.SpeechFile, nil
}

func (p *SpeechGenerationProcess) SetNext(handler Process) {
	p.next = handler
}

func (p *SpeechGenerationProcess) generateSpeech(message, folder string) (string, error) {
	filename := fmt.Sprintf("%s/%s.mp3", folder, uuid.NewString())

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
