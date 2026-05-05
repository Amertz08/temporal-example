package main

import (
	"net/http"

	"github.com/Amertz08/temporal-example/internal/api/handlers"
	"github.com/labstack/echo/v5"
)

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

	e.POST("/case", handlers.CreateCase(repo))
	e.GET("/case/:id", handlers.GetCaseById(repo))

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
