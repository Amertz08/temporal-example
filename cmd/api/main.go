package main

import (
	"log/slog"
	"net/http"

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

	var repo CaseRepository
	repo = NewInMemoryDB()

	e.POST("/case", func(c *echo.Context) error {
		// read the request body into the new Case
		var req Case
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// add the case to dbStruct
		id, err := repo.Save(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		// return case object
		return c.JSON(http.StatusOK, CaseResponse{Id: id, Case: req})
	})
	e.GET("/case/:id", func(c *echo.Context) error {
		id := c.Param("id")
		dbr, err := repo.Get(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, "not found")
		}
		return c.JSON(http.StatusOK, dbr)
	})

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
