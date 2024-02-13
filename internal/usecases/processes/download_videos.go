package processes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadVideos(topic string) (*[][]byte, error) {
	log.Println("Starting downloading videos:", topic)

	videoUrls, err := fetchVideosUrls(topic)
	if err != nil {
		return nil, fmt.Errorf("error fetching videos: %v", err)
	}

	if len(videoUrls) == 0 {
		return nil, fmt.Errorf("no videos found for topic: %s", topic)
	}

	var videos [][]byte
	for _, video := range videoUrls {
		videoBytes, err := downloadVideo(video)
		if err != nil {
			return nil, fmt.Errorf("error downloading video: %v", err)
		}

		videos = append(videos, videoBytes)
	}

	return &videos, nil
}

func fetchVideosUrls(topic string) ([]string, error) {
	pixbayVideoAPI := os.Getenv("PIXBAY_VIDEO_API")
	pixbayAPIKey := os.Getenv("PIXBAY_API_KEY")

	url := fmt.Sprintf(
		"%s/videos/?key=%s&q=%s+is&orientation=vertical&per_page=3&min_width=720&min_height=1280",
		pixbayVideoAPI,
		pixbayAPIKey,
		topic,
	)

	// Make GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch videos: %v", err)
	}
	defer resp.Body.Close()

	// Decode API response
	var videoResponse fetchVideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&videoResponse); err != nil {
		return nil, fmt.Errorf("failed to decode video response: %v", err)
	}

	// Extract video URLs from response
	var videoUrls []string
	for _, hit := range videoResponse.Hits {
		videoUrls = append(videoUrls, hit.Videos.Medium.URL)
	}

	return videoUrls, nil
}

func downloadVideo(url string) ([]byte, error) {
	// Make GET request to download the video
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download video: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body to get the video data
	video, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return video, nil
}

type hit struct {
	Videos struct {
		Duration int
		Medium   struct {
			URL string `json:"url"`
		} `json:"medium"`
	} `json:"videos"`
}

type fetchVideosResponse struct {
	Total     int   `json:"total"`
	TotalHits int   `json:"totalHits"`
	Hits      []hit `json:"hits"`
}
