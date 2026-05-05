package handlers

import (
	"log"
	"net/http"

	"github.com/Amertz08/temporal-example/internal/models"
	"github.com/labstack/echo/v5"
	"go.temporal.io/sdk/client"
)

type CaseResponse struct {
	Id string `json:"id"`
	models.Case
}

type CaseRepository interface {
	Save(models.Case) (string, error)
	Get(string) (models.Case, error)
	Close() error
}

func CreateCase(repo CaseRepository) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// read the request body into the new Case
		var req models.Case
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

		// TODO: have a unique ID here which is probably the case ID
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
	}
}
