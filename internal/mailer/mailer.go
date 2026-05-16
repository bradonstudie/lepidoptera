package mailer

import "context"

type Email struct {
	To      string
	Subject string
	Body    string
}

type Mailer interface {
	Send(ctx context.Context, email Email) error
}
