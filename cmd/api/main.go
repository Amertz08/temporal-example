package main

import (
	"log/slog"

	"github.com/labstack/echo/v5"
)

type CaseResponse struct {
	Id string `json:"id"`
	Case
}

type CaseRepository interface {
	Save(Case) (string, error)
	Get(string) (Case, error)
}

func main() {
	e := echo.New()

	repo := NewInMemoryDB()

	e.POST("/case", func(c *echo.Context) error {
		// read the request body into the new Case
		var req Case
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, err)
		}
		// add the case to dbStruct
		id, err := repo.Save(req)
		if err != nil {
			return c.JSON(500, err)
		}
		// return case object
		return c.JSON(200, CaseResponse{Id: id, Case: req})
	})
	e.GET("/case/:id", func(c *echo.Context) error {
		id := c.Param("id")
		dbr, err := repo.Get(id)
		if err != nil {
			return c.JSON(404, "not found")
		}
		return c.JSON(200, dbr)
	})

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
