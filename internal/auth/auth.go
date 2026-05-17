package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const tokenExpiry = 15 * time.Minute

// Generate 15 minute time limited HMAC token
// token format: {email}:{expiry_unix}:{signature}
func GenerateLoginToken(email, secret string) string {
	expiry := time.Now().Add(tokenExpiry).Unix()
	payload := fmt.Sprintf("%s:%d", email, expiry)
	signature := sign(payload, secret)
	return fmt.Sprintf("%s:%d:%s", email, expiry, signature)
}

func ValidateLoginToken(token, secret string) (string, error) {
	parts := strings.Split(string(token), ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	email := parts[0]
	expiry := parts[1]
	signature := parts[2]

	// verify signature
	payload := fmt.Sprintf("%s:%s", email, expiry)
	expected := sign(payload, secret)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return "", fmt.Errorf("invalid token signature")
	}

	// verify expiry status
	var expiryUnix int64
	fmt.Sscanf(expiry, "%d", &expiryUnix)
	if time.Now().Unix() > expiryUnix {
		return "", fmt.Errorf("token expired")
	}

	return email, nil
}

func sign(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
