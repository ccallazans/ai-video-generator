package processes

import (
	"log"
	"os/exec"
	"strings"
)

type LocalTextGeneration struct{}

func NewLocalTextGeneration() TextProcess {
	return &LocalTextGeneration{}
}

func (p *LocalTextGeneration) Execute(command string) (string, error) {
	log.Println("Starting text process")

	result, err := p.generateText(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *LocalTextGeneration) generateText(message string) (string, error) {

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
