package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type CaseResponse struct {
	Id string `json:"id"`
	Case
}

type UpdateRequest struct {
	Approved bool `json:"approved"`
}

func NewServer(repo CaseRepository) *echo.Echo {
	e := echo.New()

	e.POST("/case", func(c *echo.Context) error {
		// read the request body into the new Case
		var req Case
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// add the case to the database
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

	e.PUT("/case/:id", func(c *echo.Context) error {
		id := c.Param("id")
		dbr, err := repo.Get(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, "not found")
		}

		var req UpdateRequest
		if err = c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		dbr.Approved = req.Approved

		if _, err = repo.Save(dbr); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, dbr)
	})

	return e
}
