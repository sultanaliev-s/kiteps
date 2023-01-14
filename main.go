package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sultanaliev-s/kiteps/pkg/health"
)

func main() {
	e := echo.New()

	e.GET("/", echo.WrapHandler(health.NewHTTPHandler("mailer", nil)))
	e.Start(":8080")
}
