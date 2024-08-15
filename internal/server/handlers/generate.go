package handlers

import (
	"net/http"

	"github.com/ccallazans/ai-video-generator/internal/usecases"
	"github.com/labstack/echo/v4"
)

type generateRequest struct {
	Message string `json:"message"`
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

	video, err := usecases.Generate(request.Message)
	if err != nil {
		c.Logger().Errorf("Failed to generate video: %v", err)
		return respondWithError(c, http.StatusInternalServerError, "Failed to generate video")
	}

	return c.JSON(http.StatusOK, generateResponse{Video: video})
}

func respondWithError(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, map[string]string{"error": message})
}
