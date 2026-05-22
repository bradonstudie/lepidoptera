package service

import "errors"

var (
	// shows
	ErrShowNotFound = errors.New("show not found")

	// subscribers
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrAlreadySubscribed = errors.New("already subscribed")

	// admin
	ErrInvalidID     = errors.New("invalid id")
	ErrVenueNotFound = errors.New("venue not found")
	ErrBandNotFound  = errors.New("band not found")

	// shared
	ErrNotFound      = errors.New("not found")
	ErrInternalError = errors.New("internal error")
)
