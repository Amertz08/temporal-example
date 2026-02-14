package main

import (
	"log/slog"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.POST("/case", func(c *echo.Context) error {
		return c.String(200, "Case created")
	})

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
