package main

import (
	"log"
	"net/http"

	"github.com/Amertz08/temporal-example/internal/database"
	"github.com/labstack/echo/v5"
	"go.temporal.io/sdk/client"
)

type CaseResponse struct {
	Id string `json:"id"`
	database.Case
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
		var req database.Case
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// add the case to the database
		id, err := repo.Save(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		tClient, _ := client.Dial(client.Options{})
		defer tClient.Close()

		opts := client.StartWorkflowOptions{
			ID:        "case-workflow",
			TaskQueue: "test-workflow",
		}

		we, err := tClient.ExecuteWorkflow(c.Request().Context(), opts, "case-workflow", id)
		if err != nil {
			log.Println("Failed to start workflow", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

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
