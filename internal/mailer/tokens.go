package mailer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const confirmTokenExpiry = 72 * time.Hour

// confirm tokens expire after 72h; unsub tokens are permanent (embedded in sent emails).
// signatures cover all fields including expiry, so neither can be tampered with independently.
// confirm token format: confirm:{subscriberID}:{expiry_unix}:{signature}
// unsub token format:   unsub:{subscriberID}:{signature}

func ConfirmToken(subscriberID uuid.UUID, secret string) string {
	expiry := time.Now().Add(confirmTokenExpiry).Unix()
	payload := fmt.Sprintf("confirm:%s:%d", subscriberID.String(), expiry)
	sig := signToken(payload, secret)
	return fmt.Sprintf("confirm:%s:%d:%s", subscriberID.String(), expiry, sig)
}

func UnsubscribeToken(subscriberID uuid.UUID, secret string) string {
	payload := fmt.Sprintf("unsub:%s", subscriberID.String())
	sig := signToken(payload, secret)
	return fmt.Sprintf("unsub:%s:%s", subscriberID.String(), sig)
}

func ValidateConfirmTokenID(token, secret string) (uuid.UUID, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 4 {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	if parts[0] != "confirm" {
		return uuid.Nil, fmt.Errorf("invalid token prefix")
	}

	subscriberID, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subscriber id")
	}

	payload := fmt.Sprintf("confirm:%s:%s", parts[1], parts[2])
	expected := signToken(payload, secret)
	if !hmac.Equal([]byte(parts[3]), []byte(expected)) {
		return uuid.Nil, fmt.Errorf("invalid token signature")
	}

	// expiry checked after signature to avoid leaking timing info on malformed tokens
	var expiryUnix int64
	fmt.Sscanf(parts[2], "%d", &expiryUnix)
	if time.Now().Unix() > expiryUnix {
		return uuid.Nil, fmt.Errorf("token expired")
	}

	return subscriberID, nil
}

func ValidateUnsubscribeTokenID(token, secret string) (uuid.UUID, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	if parts[0] != "unsub" {
		return uuid.Nil, fmt.Errorf("invalid token prefix")
	}

	subscriberID, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subscriber id")
	}

	payload := fmt.Sprintf("unsub:%s", parts[1])
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
