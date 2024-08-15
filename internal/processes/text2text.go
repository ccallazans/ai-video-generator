package processes

import (
	"errors"
	"log"
	"os/exec"
	"strings"
)

type TextGenerationProcess struct {
	next Process
}

func NewTextGenerationProcess() *TextGenerationProcess {
	return &TextGenerationProcess{}
}

func (p *TextGenerationProcess) Execute(request interface{}) (interface{}, error) {
	context, ok := request.(*GenerationContext)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	generatedText, err := p.generateText(context.Prompt)
	if err != nil {
		return nil, err
	}
	context.Text = generatedText

	if p.next != nil {
		return p.next.Execute(context)
	}

	return context.Text, nil
}

func (p *TextGenerationProcess) SetNext(handler Process) {
	p.next = handler
}

func (p *TextGenerationProcess) generateText(message string) (string, error) {

	args := []string{
		"./pkg/llm.py",
		message,
	}

	cmd := exec.Command("python", args...)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing script local text generation: ", err)
		return "", err
	}

	response := string(cmdOutput)
	response = strings.ReplaceAll(response, "\n\n", "")

	return response, nil
}
