package processes

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	BG_VIDEO_WITH_AUDIO   = "background-video-with-audio.mp4"
	VIDEO_ASPECT_RATIO    = "aspect-ratio-video.mp4"
	VIDEO_LENGTH_CROP     = "cropped-length-video.mp4"
	VIDEO_FOLDER          = "./resources/videos"
)

type VideoGenerationProcess struct {
	next Process
}

func NewVideoGenerationProcess() *VideoGenerationProcess {
	return &VideoGenerationProcess{}
}

func (p *VideoGenerationProcess) Execute(request interface{}) (interface{}, error) {
	context, ok := request.(*GenerationContext)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	finalVideo, err := p.generateVideo(context)
	if err != nil {
		return nil, err
	}

	context.GeneratedVideo = finalVideo
	if p.next != nil {
		return p.next.Execute(context)
	}

	return context.GeneratedVideo, nil
}

func (p *VideoGenerationProcess) SetNext(handler Process) {
	p.next = handler
}

func (p *VideoGenerationProcess) generateVideo(context *GenerationContext) (string, error) {
	backgroundVideo, err := getRandomVideo(VIDEO_FOLDER)
	if err != nil {
		return "", err
	}

	bgVideoWithAudio, err := overwriteVideoAudio(context.TempDir, backgroundVideo, context.SpeechFile)
	if err != nil {
		return "", fmt.Errorf("error executing overwriteVideoAudio: %w", err)
	}

	// Crop to aspect ratio if specified
	videoToProcess := bgVideoWithAudio
	if context.AspectRatio != "" {
		aspectRatioVideo, err := cropToAspectRatio(context.TempDir, bgVideoWithAudio, context.AspectRatio)
		if err != nil {
			return "", fmt.Errorf("error executing cropToAspectRatio: %w", err)
		}
		videoToProcess = aspectRatioVideo
	}

	croppedVideo, err := cropVideoLength(context.TempDir, videoToProcess, context.SpeechFile)
	if err != nil {
		return "", fmt.Errorf("error executing cropVideoLength: %w", err)
	}

	finalPath, err := runAutocap(croppedVideo)
	if err != nil {
		return "", fmt.Errorf("error executing runAutocap: %w", err)
	}

	return finalPath, nil
}

func getRandomVideo(folder string) (string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}

	var videos []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mp4") {
			videos = append(videos, filepath.Join(folder, file.Name()))
		}
	}

	if len(videos) == 0 {
		return "", fmt.Errorf("no video files found in folder %s", folder)
	}

	rand.Seed(time.Now().Unix())
	return videos[rand.Intn(len(videos))], nil
}

func overwriteVideoAudio(tempDir, videoPath, audioPath string) (string, error) {
	videoWithAudio := filepath.Join(tempDir, BG_VIDEO_WITH_AUDIO)
	args := []string{
		"-i", videoPath,
		"-i", audioPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-map", "0:v:0",
		"-map", "1:a:0",
		videoWithAudio,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", fmt.Errorf("error overwriting video audio: %w", err)
	}

	return videoWithAudio, nil
}

func cropVideoLength(tempDir, videoPath, audioPath string) (string, error) {
	audioDuration, err := getAudioDuration(audioPath)
	if err != nil {
		return "", fmt.Errorf("error getting audio duration: %w", err)
	}

	croppedVideo := filepath.Join(tempDir, VIDEO_LENGTH_CROP)
	args := []string{
		"-i", videoPath,
		"-t", audioDuration,
		"-c", "copy",
		croppedVideo,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", fmt.Errorf("error cropping video length: %w", err)
	}

	return croppedVideo, nil
}

func getAudioDuration(audioPath string) (string, error) {
	args := []string{
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		audioPath,
	}

	output, err := executeCommandOutput("ffprobe", args)
	if err != nil {
		return "", fmt.Errorf("error getting audio duration: %w", err)
	}

	return strings.TrimSpace(output), nil
}

func runAutocap(videoPath string) (string, error) {
	finalVideo := fmt.Sprintf("%s.mp4", generateRandomWord(7))

	args := []string{
		"./scripts/captions.py",
		"attach",
		videoPath,
		finalVideo,
	}

	if err := executeCommand("python", args); err != nil {
		log.Println("Argumentos video", args)
		return "", fmt.Errorf("error running autocap: %w", err)
	}

	return finalVideo, nil
}

func executeCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command %s: %v", name, err)
		log.Printf("Command output: %s", string(output))
		return err
	}
	return nil
}

func executeCommandOutput(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command %s: %v", name, err)
		return "", err
	}
	return string(output), nil
}

func generateRandomWord(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	word := make([]byte, length)
	for i := range word {
		word[i] = letters[rand.Intn(len(letters))]
	}

	return string(word)
}

func cropToAspectRatio(tempDir, videoPath, aspectRatio string) (string, error) {
	outputVideo := filepath.Join(tempDir, VIDEO_ASPECT_RATIO)

	var width, height int
	var cropFilter string

	switch aspectRatio {
	case "16:9":
		// Standard YouTube horizontal (1920x1080)
		width, height = 1920, 1080
		// If aspect < target: scale by width, else scale by height
		cropFilter = fmt.Sprintf("scale='if(lt(a,16/9),%d,-1)':'if(lt(a,16/9),-1,%d)',crop=%d:%d", width, height, width, height)
	case "9:16":
		// Vertical for Shorts/TikTok/Reels (1080x1920)
		width, height = 1080, 1920
		cropFilter = fmt.Sprintf("scale='if(lt(a,9/16),%d,-1)':'if(lt(a,9/16),-1,%d)',crop=%d:%d", width, height, width, height)
	case "1:1":
		// Square for Instagram (1080x1080)
		width, height = 1080, 1080
		cropFilter = fmt.Sprintf("scale='if(lt(a,1),%d,-1)':'if(lt(a,1),-1,%d)',crop=%d:%d", width, height, width, height)
	case "4:3":
		// Classic 4:3 (1440x1080)
		width, height = 1440, 1080
		cropFilter = fmt.Sprintf("scale='if(lt(a,4/3),%d,-1)':'if(lt(a,4/3),-1,%d)',crop=%d:%d", width, height, width, height)
	default:
		// Default to 16:9
		width, height = 1920, 1080
		cropFilter = fmt.Sprintf("scale='if(lt(a,16/9),%d,-1)':'if(lt(a,16/9),-1,%d)',crop=%d:%d", width, height, width, height)
	}

	args := []string{
		"-i", videoPath,
		"-vf", cropFilter,
		"-c:a", "copy",
		"-y", // Overwrite output file
		outputVideo,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", fmt.Errorf("error cropping video to aspect ratio %s: %w", aspectRatio, err)
	}

	return outputVideo, nil
}
