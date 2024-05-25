package api

import (
	"net/http"

	"github.com/ccallazans/ai-video-generator/internal/generation/usecases"
	"github.com/labstack/echo/v4"
)

type generateHandlerRequest struct {
	Message string `json:"message"`
}

type generateHandlerResponse struct {
	Video string `json:"video"`
}

func GenerateHandler(c echo.Context) error {
	var request generateHandlerRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	if request.Message == "" {
		return c.JSON(http.StatusBadRequest, "Message is a required field")
	}

	video, err := usecases.Generate(request.Message)
	if err != nil {
		c.Logger().Errorf("Failed to generate video: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to generate video")
	}

	return c.JSON(http.StatusOK, generateHandlerResponse{
		Video: video,
	})
}
