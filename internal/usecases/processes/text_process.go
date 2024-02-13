package processes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func TextProcess(message string) (string, error) {
	log.Println("Starting text process:", message)

	text, err := getText(message)
	if err != nil {
		return "", fmt.Errorf("error processing text: %v", err)
	}

	return text, nil
}

func getText(message string) (string, error) {
	textGenerationAPI := os.Getenv("TEXT_GENERATION_API")

	// Marshal request body
	requestBody, err := json.Marshal(map[string]string{"text": message})
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON request: %v", err)
	}

	// Make POST request to the text generation API
	resp, err := http.Post(textGenerationAPI, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response body
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("error decoding response body: %v", err)
	}

	// Extract generated text from response
	generatedText, ok := responseBody["response"]
	if !ok {
		return "", fmt.Errorf("response field missing in API response")
	}

	return generatedText, nil
}
