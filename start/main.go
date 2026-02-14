package main

import (
	"context"
	"log"

	"github.com/Amertz08/temporal-example/greeting"
	"go.temporal.io/sdk/client"
)

// This is what is actually kicking off a workflow run
func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client:", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "greating-workflow",
		TaskQueue: "my-task-queue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, greeting.SayHelloWorkFlow, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow:", err)
	}
	log.Println("Workflow execution started:", we.GetID(), "RunID: ", we.GetRunID())

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get workflow result:", err)
	}
	log.Println("Workflow result:", result)
}
