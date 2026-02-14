package temporal_example

import "context"

func Greet(ctx context.Context) (string, error) {
	return "Hello World!", nil
}
