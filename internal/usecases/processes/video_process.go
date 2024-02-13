package processes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ccallazans/ai-video-generator/internal/utils"
)

func VideoProcess(videosBytes *[][]byte, speechBytes *[]byte) (string, error) {
	log.Println("Starting video process")

	// Generate a temporary directory for storing intermediate files.
	tempDir, err := generateTemporaryFolder()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer os.RemoveAll(tempDir)

	// Generate temporary files for videos and speech.
	videosTextFilename, speechFilename, err := generateTemporaryFiles(videosBytes, speechBytes, tempDir)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Merge videos into a single file.
	mergedVideosFilename, err := mergeVideos(videosTextFilename, tempDir)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Add speech to the merged video.
	videoWithSpeechFilename, err := addSpeechToVideo(mergedVideosFilename, speechFilename, tempDir)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Get the duration of the speech file.
	speechDuration, err := getSpeechTimeDuration(speechFilename)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Crop the video to match the duration of the speech.
	video, err := cropVideoDuration(videoWithSpeechFilename, speechDuration, tempDir)
	if err != nil {
		log.Println(err)
		return "", err
	}

	uploadFile(video)

	return video, nil
}

func generateTemporaryFolder() (string, error) {
	log.Println("Generating temporary folder")

	dirName, err := os.MkdirTemp("", "merge-videos")
	if err != nil {
		log.Printf("Failed to create temporary directory: %v\n", err)
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	return dirName, nil
}

func generateTemporaryFiles(videosBytes *[][]byte, speechBytes *[]byte, dirName string) (string, string, error) {
	log.Println("Generating temporary files")

	var videoFilenames []string
	for _, videoData := range *videosBytes {
		videoFile, err := os.CreateTemp(dirName, "*.mp4")
		if err != nil {
			log.Fatalf("Error creating temporary video file: %v", err)
			return "", "", err
		}
		defer videoFile.Close()

		if _, err := videoFile.Write(videoData); err != nil {
			log.Fatalf("Error writing video data to temporary file: %v", err)
			return "", "", err
		}

		newName := strings.Split(videoFile.Name(), ".mp4")[0] + "crop.mp4"
		args := []string{
			"-i", videoFile.Name(),
			"-vf", "crop=720:1280:in_w/2-360:in_h/2-640",
			"-c:a", "copy",
			newName,
		}

		cmd := exec.Command("ffmpeg", args...)
		_, err = cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
			return "", "", err
		}

		videoFilenames = append(videoFilenames, newName)
	}

	speechFile, err := os.CreateTemp(dirName, "*.mp3")
	if err != nil {
		log.Fatalf("Error creating temporary speech file: %v", err)
		return "", "", err
	}
	defer speechFile.Close()

	if _, err := speechFile.Write(*speechBytes); err != nil {
		log.Fatalf("Error writing speech data to temporary file: %v", err)
		return "", "", err
	}

	textFile, err := os.CreateTemp(dirName, "*.txt")
	if err != nil {
		log.Fatalf("Error creating temporary text file: %v", err)
		return "", "", err
	}
	defer textFile.Close()

	textData := ""
	for _, vid := range videoFilenames {
		textData += fmt.Sprintf("file %s\n", vid)
	}

	if _, err := textFile.Write([]byte(textData)); err != nil {
		log.Fatalf("Error writing text data to temporary file: %v", err)
		return "", "", err
	}

	return textFile.Name(), speechFile.Name(), nil
}

func mergeVideos(file string, dirName string) (string, error) {
	log.Println("Merging videos")

	mergedFile := fmt.Sprintf("%s/%s.mp4", dirName, utils.RandomString())

	args := []string{
		"-f", "concat",
		"-safe", "0",
		// "-map", "0:v",
		"-i", file,
		"-c", "copy",
		mergedFile,
	}

	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return mergedFile, nil
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
	log.Println(args)
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		return "", err
	}

	return finalVideo, nil
}

func uploadFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new buffer to store the file contents
	fileBuffer := bytes.Buffer{}
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		return "", err
	}

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the multipart form
	fileWriter, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fileWriter, &fileBuffer)
	if err != nil {
		return "", err
	}

	// Close the multipart writer
	writer.Close()

	// Create the request
	req, err := http.NewRequest("POST", "http://127.0.0.1:5000/add-subtitle", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create a new file to save the response
	responseFile, err := os.Create(fmt.Sprintf("./generated/%s.mp4", utils.RandomString()))
	if err != nil {
		return "", err
	}
	defer responseFile.Close()

	// Copy the response body to the file
	_, err = io.Copy(responseFile, resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("File uploaded successfully and response saved as %s\n", responseFile.Name())
	return responseFile.Name(), nil
}
