package processes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func SpeechProcess(message string) (*[]byte, error) {
	log.Println("Starting speech process:", message)

	audioBytes, err := getSpeech(message)
	if err != nil {
		return nil, fmt.Errorf("error processing speech: %v", err)
	}

	return &audioBytes, nil
}

func getSpeech(message string) ([]byte, error) {
	speechGenerationAPI := os.Getenv("SPEECH_GENERATION_API")

	// Marshal request body
	requestBody, err := json.Marshal(map[string]string{"text": message})
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON request: %v", err)
	}

	// Make POST request to the speech generation API
	resp, err := http.Post(speechGenerationAPI, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read response body to get speech audio bytes
	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return audio, nil
}
