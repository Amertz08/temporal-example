package activities

import (
	"context"
	"log"
	"time"
)

func SendEmail(ctx context.Context, to, subject, body string) error {
	log.Printf("sending email to %s subject %s", to, subject)
	time.Sleep(2 * time.Second)
	return nil
}
