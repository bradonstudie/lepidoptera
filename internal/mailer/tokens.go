package mailer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func ConfirmToken(subscriberID uuid.UUID, secret string) string {
	return signToken("confirm:"+subscriberID.String(), secret)
}

func UnsubscribeToken(subscriberID uuid.UUID, secret string) string {
	return signToken("unsub:"+subscriberID.String(), secret)
}

func signToken(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func ValidateToken(prefix string, token string, secret string, subscriberID uuid.UUID) bool {
	expected := signToken(prefix+":"+subscriberID.String(), secret)
	return hmac.Equal([]byte(token), []byte(expected))
}
