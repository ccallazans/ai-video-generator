package processes

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	BG_VIDEO_PATH  = "./resources/videos"
	GENERATED_PATH = "./generated"
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
	bgVideo, err := selectBackgroundVideo()
	if err != nil {
		return "", err
	}

	videoWithSpeech, err := addSpeechToVideo(bgVideo, context.SpeechFile, context.TempDir)
	if err != nil {
		return "", err
	}

	speechDuration, err := getSpeechTimeDuration(context.SpeechFile)
	if err != nil {
		return "", err
	}

	croppedVideo, err := cropVideoDuration(videoWithSpeech, speechDuration, context.TempDir)
	if err != nil {
		return "", err
	}

	captionsVideo, err := addVideoSubtitles(croppedVideo, context.TempDir)
	if err != nil {
		return "", err
	}

	finalVideo, err := exportMp4(captionsVideo)
	if err != nil {
		return "", err
	}

	return finalVideo, nil
}

func selectBackgroundVideo() (string, error) {
	log.Println("Selecting background video")

	bgVideos, err := os.ReadDir(BG_VIDEO_PATH)
	if err != nil {
		return "", fmt.Errorf("error reading background videos directory: %w", err)
	}

	if len(bgVideos) == 0 {
		return "", errors.New("no files found in the background videos directory")
	}

	randomIndex := rand.Intn(len(bgVideos))
	bgFilename := bgVideos[randomIndex].Name()
	bgPath := filepath.Join(BG_VIDEO_PATH, bgFilename)

	return bgPath, nil
}

func addSpeechToVideo(videoFile, speechFile, tempDir string) (string, error) {
	log.Println("Adding speech to video")

	videoWithSpeechFile := filepath.Join(tempDir, uuid.NewString()+".mp4")

	args := []string{
		"-i", videoFile,
		"-i", speechFile,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "experimental",
		"-map", "0:v:0",
		"-map", "1:a:0",
		videoWithSpeechFile,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", err
	}

	return videoWithSpeechFile, nil
}

func getSpeechTimeDuration(speechFile string) (string, error) {
	log.Println("Getting speech time duration")

	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		speechFile,
	}

	output, err := executeCommandOutput("ffprobe", args)
	if err != nil {
		return "", err
	}

	durationSeconds, err := strconv.ParseFloat(strings.TrimSpace(output), 64)
	if err != nil {
		return "", err
	}

	hours := int(math.Floor(durationSeconds / 3600))
	minutes := int(math.Floor((durationSeconds - float64(hours)*3600) / 60))
	remainingSeconds := durationSeconds - float64(hours)*3600 - float64(minutes)*60
	milliseconds := int((remainingSeconds - math.Floor(remainingSeconds)) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, int(remainingSeconds), milliseconds), nil
}

func cropVideoDuration(videoFile, duration, tempDir string) (string, error) {
	log.Println("Cropping video duration")

	finalVideo := filepath.Join(tempDir, uuid.NewString()+".mp4")

	args := []string{
		"-i", videoFile,
		"-c", "copy",
		"-t", duration,
		finalVideo,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", err
	}

	return finalVideo, nil
}

func addVideoSubtitles(videoFile, tempDir string) (string, error) {
	log.Println("Adding subtitles")

	captionsVideo := filepath.Join(tempDir, uuid.NewString()+".mp4")

	args := []string{
		"./pkg/captions.py",
		videoFile,
		captionsVideo,
	}

	if err := executeCommand("python", args); err != nil {
		return "", err
	}

	return captionsVideo, nil
}

func exportMp4(videoFile string) (string, error) {
	log.Println("Exporting to mp4")

	finalVideoPath := filepath.Join(GENERATED_PATH, uuid.NewString()+".mp4")

	args := []string{
		"-i", videoFile,
		"-codec", "copy",
		finalVideoPath,
	}

	if err := executeCommand("ffmpeg", args); err != nil {
		return "", err
	}

	return finalVideoPath, nil
}

func executeCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Error executing command %s: %v", name, err)
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
