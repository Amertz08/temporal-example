package main

import (
	"github.com/Amertz08/temporal-example/internal/api/handlers"
	"github.com/labstack/echo/v5"
)

func NewServer(repo CaseRepository) *echo.Echo {
	e := echo.New()

	e.POST("/case", handlers.CreateCase(repo))
	e.GET("/case/:id", handlers.GetCaseById(repo))
	e.PATCH("/case/:id", handlers.UpdateCase(repo))

	return e
}
