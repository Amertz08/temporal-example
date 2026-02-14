package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type CaseResponse struct {
	Id string `json:"id"`
	Case
}

type PatchRequest struct {
	Name         *string `json:"name"`
	Address      *string `json:"address"`
	Email        *string `json:"email"`
	VinNumber    *string `json:"vin_number"`
	Approved     *bool   `json:"approved"`
	Manufactured *bool   `json:"manufactured"`
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

	e.PATCH("/case/:id", func(c *echo.Context) error {
		id := c.Param("id")
		dbr, err := repo.Get(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, "not found")
		}

		var req PatchRequest
		if err = c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// Only update fields that were provided
		if req.Name != nil {
			dbr.Name = *req.Name
		}
		if req.Address != nil {
			dbr.Address = *req.Address
		}
		if req.Email != nil {
			dbr.Email = *req.Email
		}
		if req.VinNumber != nil {
			dbr.VinNumber = *req.VinNumber
		}
		if req.Approved != nil {
			dbr.Approved = *req.Approved
		}
		if req.Manufactured != nil {
			dbr.Manufactured = *req.Manufactured
		}

		if _, err = repo.Save(dbr); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, dbr)
	})

	return e
}
