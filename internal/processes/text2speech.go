package processes

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"

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

	speechFilename, err := p.generateSpeech(context.Text, context.TempDir, context.Voice)
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

func (p *SpeechGenerationProcess) generateSpeech(text, folder, voice string) (string, error) {
	filename := filepath.Join(folder, uuid.NewString()+".mp3")

	args := []string{
		"./scripts/tts.py",
		text,
		filename,
		voice,
	}

	cmd := exec.Command("python", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error text2speech: ", args)
		log.Println("Command output: ", string(output))
		return "", err
	}

	return filename, nil
}
