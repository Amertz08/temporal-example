package greeting

import "context"

func Greet(ctx context.Context) (string, error) {
	return "Hello World!", nil
}
