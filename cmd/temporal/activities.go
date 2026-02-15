package main

import (
	"context"
	"log"
)

func SendEmail(ctx context.Context, to, subject, body string) error {
	log.Printf("sending email to %s subject %s", to, subject)
	return nil
}
