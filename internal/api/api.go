package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApi() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	setupV1Routes(e.Group("/api/v1"))

	return e
}

func setupV1Routes(v1 *echo.Group) {
	v1.POST("/generate", GenerateHandler)
}
