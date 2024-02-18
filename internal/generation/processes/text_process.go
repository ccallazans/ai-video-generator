package processes

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

	requestData := request{
		Model:  "llama2",
		Prompt: message,
		Stream: false,
	}

	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		log.Println("Error marshalling request data:", err)
		return "", err
	}

	url := "http://ollama:11434/api/generate"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	var responseData response
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		log.Println("Error parsing response JSON:", err)
		return "", err
	}

	return responseData.Response, nil
}

type request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type response struct {
	Response string `json:"response"`
}
