package service

import (
	"context"
	"errors"
	"fmt"

	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
)

var ErrInvalidToken = errors.New("invalid or expired token")
var ErrAlreadySubscribed = errors.New("already subscribed")

type SubscriberService struct {
	queries *db.Queries
	mailer  mailer.Mailer
	secret  string
}

func NewSubscriberService(queries *db.Queries, m mailer.Mailer, secret string) *SubscriberService {
	return &SubscriberService{queries: queries, mailer: m, secret: secret}
}

func (s *SubscriberService) Subscribe(ctx context.Context, email string) error {
	sub, err := s.queries.CreateSubscriber(ctx, email)
	if err != nil {
		return fmt.Errorf("Subscribe: %w", err)
	}

	// zero UUID means ON CONFLICT DO NOTHING fired — already subscribed
	if sub.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return ErrAlreadySubscribed
	}

	token := mailer.ConfirmToken(sub.ID, s.secret)
	confirmEmail := mailer.ConfirmSubscriptionEmail(token)
	confirmEmail.To = email
	s.mailer.Send(ctx, confirmEmail)

	return nil
}

func (s *SubscriberService) Confirm(ctx context.Context, token string) error {
	id, err := mailer.ValidateConfirmTokenID(token, s.secret)
	if err != nil {
		return ErrInvalidToken
	}

	if err := s.queries.ConfirmSubscriber(ctx, id); err != nil {
		return fmt.Errorf("Confirm: %w", err)
	}

	sub, err := s.queries.GetSubscriberByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Confirm get subscriber: %w", err)
	}

	unsubToken := mailer.UnsubscribeToken(sub.ID, s.secret)
	confirmedEmail := mailer.SubscriptionConfirmedEmail(unsubToken)
	confirmedEmail.To = sub.Email
	s.mailer.Send(ctx, confirmedEmail)

	return nil
}

func (s *SubscriberService) Unsubscribe(ctx context.Context, token string) error {
	id, err := mailer.ValidateUnsubscribeTokenID(token, s.secret)
	if err != nil {
		return ErrInvalidToken
	}

	if err := s.queries.UnsubscribeSubscriber(ctx, id); err != nil {
		return fmt.Errorf("Unsubscribe: %w", err)
	}

	return nil
}
