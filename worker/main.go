package main

import (
	"log"

	"github.com/Amertz08/temporal-example/greeting"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// The worker is responsible for executing workflows and activities
func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client:", err)
	}
	defer c.Close()

	w := worker.New(c, "my-task-queue", worker.Options{})

	w.RegisterWorkflow(greeting.SayHelloWorkFlow)
	w.RegisterActivity(greeting.Greet)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker:", err)
	}
}
