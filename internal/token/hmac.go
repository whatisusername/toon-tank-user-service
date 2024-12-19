package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func GenerateBase64HMAC(key, message string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to write HMAC: %w", err)
	}

	hmacBytes := h.Sum(nil)
	base64Encoded := base64.StdEncoding.EncodeToString(hmacBytes)

	return base64Encoded, nil
}
