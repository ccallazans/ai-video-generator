package processes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type TextGenerationProcess struct {
	next   Process
	config Config
}

type Config struct {
	OllamaURL   string
	OllamaModel string
}

func NewTextGenerationProcess() *TextGenerationProcess {
	return &TextGenerationProcess{
		config: Config{
			OllamaURL:   os.Getenv("OLLAMA_URL"),
			OllamaModel: os.Getenv("OLLAMA_MODEL"),
		},
	}
}

func (p *TextGenerationProcess) Execute(request interface{}) (interface{}, error) {
	context, ok := request.(*GenerationContext)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	// Extract language from voice and modify prompt
	language := getLanguageFromVoice(context.Voice)
	enhancedPrompt := enhancePromptWithLanguage(context.Prompt, language)

	generatedText, err := p.generateText(enhancedPrompt)
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

func (p *TextGenerationProcess) generateText(prompt string) (string, error) {
	payload := requestPayload{
		Model:  p.config.OllamaModel,
		Prompt: prompt,
		Stream: false,
	}

	var result response
	if err := makeOllamaRequest("POST", p.config.OllamaURL, payload, &result); err != nil {
		return "", fmt.Errorf("error generating text: %w", err)
	}

	return result.Response, nil
}

type requestPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type response struct {
	Response string `json:"response"`
}

func makeOllamaRequest(method, url string, payload interface{}, result interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response body: %w", err)
	}

	return nil
}

// getLanguageFromVoice extracts the language name from Edge TTS voice code
// Examples: "es-ES-ElviraNeural" -> "Spanish", "en-US-AriaNeural" -> "English"
func getLanguageFromVoice(voice string) string {
	if voice == "" {
		return "English"
	}

	// Extract language code (e.g., "es-ES" from "es-ES-ElviraNeural")
	parts := strings.Split(voice, "-")
	if len(parts) < 2 {
		return "English"
	}

	languageCode := strings.ToLower(parts[0])

	// Map language codes to full names
	languageMap := map[string]string{
		"en": "English",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
		"it": "Italian",
		"pt": "Portuguese",
		"nl": "Dutch",
		"ru": "Russian",
		"ja": "Japanese",
		"ko": "Korean",
		"zh": "Chinese",
		"ar": "Arabic",
		"hi": "Hindi",
		"tr": "Turkish",
		"pl": "Polish",
		"sv": "Swedish",
		"no": "Norwegian",
		"da": "Danish",
		"fi": "Finnish",
	}

	if language, ok := languageMap[languageCode]; ok {
		return language
	}

	return "English"
}

// enhancePromptWithLanguage adds language instruction to the prompt
func enhancePromptWithLanguage(prompt string, language string) string {
	if language == "English" {
		return prompt
	}

	return fmt.Sprintf("%s\n\nIMPORTANT: You must respond entirely in %s. Do not use English or any other language.", prompt, language)
}