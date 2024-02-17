package processes

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ccallazans/ai-video-generator/internal/utils"
)

type LocalVideoGeneration struct {
	tempFolder     string
	speechFilename string
}

func NewLocalVideoGeneration(tempFolder string, speechFilename string) VideoProcess {
	return &LocalVideoGeneration{tempFolder: tempFolder, speechFilename: speechFilename}
}

func (p *LocalVideoGeneration) Execute(command string) (string, error) {
	log.Println("Starting video process")

	result, err := p.generateVideo(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *LocalVideoGeneration) generateVideo(message string) (string, error) {
	bgVideo, err := selectBackgroundVideo()
	if err != nil {
		return "", err
	}

	videoWithSpeechFilename, err := addSpeechToVideo(bgVideo, p.speechFilename, p.tempFolder)
	if err != nil {
		log.Println(err)
		return "", err
	}

	speechDuration, err := getSpeechTimeDuration(p.speechFilename)
	if err != nil {
		log.Println(err)
		return "", err
	}

	video, err := cropVideoDuration(videoWithSpeechFilename, speechDuration, p.tempFolder)
	if err != nil {
		log.Println(err)
		return "", err
	}

	captionsVideo, err := addVideoSubtitles(video, p.tempFolder)
	if err != nil {
		log.Println(err)
		return "", err
	}

	finalVideo, err := exportMp4(captionsVideo)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return finalVideo, nil
}

func selectBackgroundVideo() (string, error) {
	log.Println("Selecting background video")

	bgVideoPath := "./resources/videos"

	bgVideos, err := os.ReadDir(bgVideoPath)
	if err != nil {
		log.Println("Error reading background videos directory:", err)
		return "", err
	}

	if len(bgVideos) == 0 {
		log.Println("No files found in the background videos directoryy")
		return "", err
	}

	randomIndex := rand.Intn(len(bgVideos))

	bgFilename := bgVideos[randomIndex].Name()
	bgPath := filepath.Join(bgVideoPath, bgFilename)

	return bgPath, nil
}

func addSpeechToVideo(mergedVideoFile string, speechFile string, dirName string) (string, error) {
	log.Println("Adding speech to video")

	videoWithSpeechFile := fmt.Sprintf("%s/%s.mp4", dirName, utils.RandomString())

	args := []string{
		"-i", mergedVideoFile,
		"-i", speechFile,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "experimental",
		"-map", "0:v:0",
		"-map", "1:a:0",
		videoWithSpeechFile,
	}

	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return videoWithSpeechFile, nil
}

func getSpeechTimeDuration(speechFilename string) (string, error) {
	log.Println("Getting speech time duration")

	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		speechFilename,
	}

	cmd := exec.Command("ffprobe", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return "", err
	}

	strOutput := strings.TrimSpace(string(output))

	durationSeconds, err := strconv.ParseFloat(strOutput, 64)
	if err != nil {
		log.Println(err)
		return "", err
	}

	hours := int(math.Floor(durationSeconds / 3600))
	minutes := int(math.Floor((durationSeconds - float64(hours)*3600) / 60))
	remainingSeconds := durationSeconds - float64(hours)*3600 - float64(minutes)*60
	milliseconds := int((remainingSeconds - math.Floor(remainingSeconds)) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, int(remainingSeconds), milliseconds), nil
}

func cropVideoDuration(video string, duration string, dirName string) (string, error) {
	log.Println("Cropping video duration")

	finalVideo := fmt.Sprintf("%s/%s.mp4", dirName, utils.RandomString())

	args := []string{
		"-i", video,
		"-c", "copy",
		"-t", duration,
		finalVideo,
	}

	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		return "", err
	}

	return finalVideo, nil
}

func addVideoSubtitles(filePath string, tempDir string) (string, error) {
	log.Println("Adding subtitles")

	filename := fmt.Sprintf("%s.mp4", utils.RandomString())

	args := []string{
		"./pkg/captions.py",
		filePath,
		filename,
	}

	cmd := exec.Command("python", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error executing script local video generation: ", err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", tempDir, filename), nil
}

func exportMp4(filePath string) (string, error) {
	log.Println("Exporting to mp4")

	finalVideoPath := fmt.Sprintf("./generated/%s.mp4", utils.RandomString())

	args := []string{
		"-i",
		filePath,
		"-codec",
		"copy",
		finalVideoPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error exporting video to mp4: ", err)
		return "", err
	}

	return finalVideoPath, nil
}
