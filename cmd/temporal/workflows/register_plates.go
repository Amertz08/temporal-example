package workflows

import (
	"fmt"
	"log"

	"github.com/Amertz08/temporal-example/cmd/temporal/activities"
	"github.com/Amertz08/temporal-example/internal/models"
	"go.temporal.io/sdk/workflow"
)

const ApprovedSignal = "approved"

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

	err = workflow.ExecuteActivity(
		ctx,
		activities.SendEmail,
		caseRecord.Email,
		"Your case has been submitted",
		"We have received your information and will be validating the information and will get back to you with further steps.",
	).Get(ctx, nil)
	if err != nil {
		fmt.Println("errored on sending initial email")
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

	// Block until internal user verifies it's fine
	workflow.GetSignalChannel(ctx, ApprovedSignal).Receive(ctx, nil)
	// If not fine do remediation workflow

	// Send email notifying user to schedule an appointment

	// Send email appointment confirmation
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

	// employee validates gov issued ID and registration
	// ID & registration uploaded
	// employee accepts payment
	// mfg order created
	// mfg started
	// mfg completed
	// shipping started
	// plate shipped
	// done

	// TODO: wait for approval
	return nil
}
