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
	log.Println("1 Starting speech process")

	result, err := p.generateSpeech(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *LocalSpeechGeneration) generateSpeech(message string) (string, error) {
	filename := fmt.Sprintf("%s/%s.wav", p.tempFolder, utils.RandomString())

	args := []string{
		"--text",
		message,
		"--model_name",
		"tts_models/en/vctk/vits",
		"--out_path",
		filename,
		"--speaker_idx",
		"p266",
	}

	cmd := exec.Command("tts", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing script local speech generation: ", err)
		return "", err
	}

	return filename, nil
}
