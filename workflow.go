package temporal_example

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func SayHelloWorkFlow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, Greet, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}
