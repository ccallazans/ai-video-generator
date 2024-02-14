package api

import (
	"net/http"

	"github.com/ccallazans/ai-video-generator/internal/usecases"
	"github.com/labstack/echo/v4"
)

type GenerateHandlerRequest struct {
	Message string `json:"message"`
}

func GenerateHandler(c echo.Context) error {
	var request GenerateHandlerRequest
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if request.Message == "" {
		return c.String(http.StatusBadRequest, "message and topic are required fields")
	}

	video, err := usecases.Generate(request.Message)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "failed to generate video")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"video": video,
	})
}
