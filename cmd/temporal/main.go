package main

import (
	"log"

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

	w.RegisterWorkflow(RegisterLicensePlateWorkflow)
	w.RegisterActivity(SendEmail)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker failed: ", err)
	}
}
