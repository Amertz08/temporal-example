package workflows

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Amertz08/temporal-example/cmd/temporal/activities"
	"github.com/Amertz08/temporal-example/internal/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	ApprovedSignal             = "approved"
	AppointmentScheduledSignal = "appointment_scheduled"
	ValidatedIdSignal          = "id_validated"
)

type RegisterLicensePlateInput struct {
	caseId string
}

func RegisterLicensePlateWorkflow(ctx workflow.Context, input RegisterLicensePlateInput) error {
	sendEmailRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumAttempts:    3,
	}
	sendEmailActivityOptions := workflow.ActivityOptions{
		RetryPolicy: sendEmailRetryPolicy,
	}

	var caseRecord *models.Case
	err := workflow.ExecuteActivity(
		ctx,
		activities.GetCaseById,
		input.caseId,
	).Get(ctx, &caseRecord)

	if err != nil {
		log.Println("Failed to get case from database", err)
		return err
	}

	ctx = workflow.WithActivityOptions(ctx, sendEmailActivityOptions)
	err = workflow.ExecuteActivity(
		ctx,
		SendInitialEmail,
		caseRecord.Email,
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

	// Calculate fee amount
	var feeAmount int64
	workflow.ExecuteActivity(
		ctx,
		activities.CalculateFeeAmount,
		&vinDetails,
	).Get(ctx, &feeAmount)

	// Send email notifying user to schedule an appointment

	// Block until appointment is schedule
	workflow.GetSignalChannel(ctx, AppointmentScheduledSignal).Receive(ctx, nil)

	// Send email appointment confirmation
	ctx = workflow.WithActivityOptions(ctx, sendEmailActivityOptions)
	err = workflow.ExecuteActivity(
		ctx,
		SendAppointmentConfirmationEmail,
		caseRecord.Email,
	).Get(ctx, nil)
	if err != nil {
		log.Println("Failed to send email", err)
		return err
	}

	// Block until employee validates gov issued ID and registration
	workflow.GetSignalChannel(ctx, ValidatedIdSignal).Receive(ctx, nil)
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

func SendInitialEmail(ctx context.Context, toEmail string) error {
	return activities.SendEmail(
		ctx,
		toEmail,
		"Your case has been submitted",
		"We have received your information and will be validating the information and will get back to you with further steps.",
	)
}

func SendAppointmentConfirmationEmail(ctx context.Context, toEmail string) error {
	return activities.SendEmail(
		ctx,
		toEmail,
		"License Plate Registered",
		"Your appointment is set for 2025-01-01 at 9:00 AM CST",
	)
}
