package handlers

import (
	"net/http"

	"github.com/ccallazans/ai-video-generator/internal/usecases"
	"github.com/labstack/echo/v4"
)

type generateRequest struct {
	Message     string `json:"message"`
	AspectRatio string `json:"aspect_ratio"` // "16:9" (default), "9:16", "1:1"
	Voice       string `json:"voice"`        // Edge TTS voice (e.g., "en-US-AriaNeural", "es-ES-ElviraNeural")
}

type generateResponse struct {
	Video string `json:"video"`
}

func GenerateHandler(c echo.Context) error {
	var request generateRequest

	if err := c.Bind(&request); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid request payload")
	}

	if request.Message == "" {
		return respondWithError(c, http.StatusBadRequest, "Message is a required field")
	}

	// Set defaults
	if request.AspectRatio == "" {
		request.AspectRatio = "16:9"
	}
	if request.Voice == "" {
		request.Voice = "en-US-AriaNeural"
	}

	video, err := usecases.Generate(request.Message, request.AspectRatio, request.Voice)
	if err != nil {
		c.Logger().Errorf("Failed to generate video: %v", err)
		return respondWithError(c, http.StatusInternalServerError, "Failed to generate video")
	}

	return c.JSON(http.StatusOK, generateResponse{Video: video})
}

func respondWithError(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, map[string]string{"error": message})
}
