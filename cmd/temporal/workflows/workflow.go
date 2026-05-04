package workflows

import (
	"fmt"
	"log"

	"github.com/Amertz08/temporal-example/cmd/temporal/activities"
	"github.com/Amertz08/temporal-example/internal/models"
	"go.temporal.io/sdk/workflow"
)

func RegisterLicensePlateWorkflow(ctx workflow.Context, caseId string) error {
	var caseRecord *models.Case
	err := workflow.ExecuteActivity(
		ctx,
		activities.GetCaseById,
		caseId,
	).Get(ctx, &caseRecord)

	if err != nil {
		log.Println("Failed to get case from database", err)
		return err
	}

	var vinDetails *models.VinDetails
	err = workflow.ExecuteActivity(
		ctx,
		activities.GetVinDetails,
		caseRecord.VinNumber,
	).Get(ctx, &vinDetails)
	if err != nil {
		fmt.Println("error getting vin details")
		return err
	}

	err = workflow.ExecuteActivity(
		ctx,
		activities.SendEmail,
		caseRecord.Email,
		"License Plate Registered",
		"Your appointment is set for 2025-01-01 at 9:00 AM CST",
	).Get(ctx, nil)
	if err != nil {
		log.Println("Failed to send email", err)
		return err
	}

	// TODO: wait for approval
	return nil
}
