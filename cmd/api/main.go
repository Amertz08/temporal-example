package main

import (
	"log/slog"

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

	// TODO: make persistent with JSON file
	db := dbStruct{}

	e.POST("/case", func(c *echo.Context) error {
		// read the request body into new Case
		var req Case
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, err)
		}
		// add the case to dbStruct
		id := "abc-def"
		db[id] = req
		// return case object

		return c.JSON(200, CaseResponse{Id: id, Case: req})
	})
	e.GET("/case/:id", func(c *echo.Context) error {
		id := c.Param("id")
		dbr, ok := db[id]
		if !ok {
			return c.JSON(404, "not found")
		}
		return c.JSON(200, dbr)
	})

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
