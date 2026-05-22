package mailer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// token format: {prefix}:{subscriberID}:{signature}

func ConfirmToken(subscriberID uuid.UUID, secret string) string {
	return buildToken("confirm", subscriberID, secret)
}

func UnsubscribeToken(subscriberID uuid.UUID, secret string) string {
	return buildToken("unsub", subscriberID, secret)
}

func buildToken(prefix string, subscriberID uuid.UUID, secret string) string {
	payload := fmt.Sprintf("%s:%s", prefix, subscriberID.String())
	sig := signToken(payload, secret)
	return fmt.Sprintf("%s:%s:%s", prefix, subscriberID.String(), sig)
}

func ValidateConfirmTokenID(token, secret string) (uuid.UUID, error) {
	return validateToken("confirm", token, secret)
}

func ValidateUnsubscribeTokenID(token, secret string) (uuid.UUID, error) {
	return validateToken("unsub", token, secret)
}

func validateToken(prefix, token, secret string) (uuid.UUID, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	if parts[0] != prefix {
		return uuid.Nil, fmt.Errorf("invalid token prefix")
	}

	subscriberID, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subscriber id")
	}

	payload := fmt.Sprintf("%s:%s", prefix, parts[1])
	expected := signToken(payload, secret)
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return uuid.Nil, fmt.Errorf("invalid token signature")
	}

	return subscriberID, nil
}

func signToken(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
