package main

import (
	"log"

	"github.com/Amertz08/temporal-example/cmd/temporal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	w := worker.New(c, "test-workflow", worker.Options{})

	w.RegisterWorkflow(workflows.RegisterLicensePlateWorkflow)
	w.RegisterActivity(workflows.SendInitialEmail)
	w.RegisterActivity(workflows.SendAppointmentConfirmationEmail)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker failed: ", err)
	}
}
