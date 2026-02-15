package main

import (
	"log"

	"github.com/Amertz08/temporal-example/internal/database"
	"go.temporal.io/sdk/workflow"
)

func RegisterLicensePlateWorkflow(ctx workflow.Context, caseId string) error {
	repo, err := database.NewJSONFileDB("../api/cases.json")
	if err != nil {
		log.Println("Failed to initialize database", err)
		return err
	}
	defer repo.Close()
	// TODO: I don't think you're supposed to query the DB in the workflow but instead do it as an activity.
	caseRecord, err := repo.Get(caseId)
	if err != nil {
		log.Println("Failed to get case from database", err)
		return err
	}
	err = workflow.ExecuteActivity(ctx, SendEmail, caseRecord.Email, "License Plate Registered", "Your appointment is set for 2025-01-01 at 9:00 AM CST").Get(ctx, nil)
	if err != nil {
		log.Println("Failed to send email", err)
		return err
	}
	return nil
}
