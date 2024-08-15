package server

import (
	"net/http"

	"github.com/ccallazans/ai-video-generator/internal/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", s.healthHandler)

	v1 := e.Group("/api/v1")
	v1.POST("/generate", handlers.GenerateHandler)

	return e
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &map[string]string{
		"Status": "OK",
	})
}
