package main

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Case struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	VinNumber string `json:"vin_number"`
}

type CaseResponse struct {
	Id string `json:"id"`
	Case
}

type dbStruct map[string]Case

func main() {
	e := echo.New()

	db := dbStruct{}

	e.POST("/case", func(c *echo.Context) error {
		// read the request body into new Case
		var req Case
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, err)
		}
		// add the case to dbStruct
		id := uuid.New().String()
		db[id] = req
		// return case object

		return c.JSON(200, CaseResponse{Id: id, Case: req})
	})

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
