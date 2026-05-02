package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func main() {
	key := randomKey()

	token := sign(key, "user_id=user_123")
	fmt.Println("token:", token)

	payload, ok := verify(key, token)
	fmt.Println("verify ok:", ok, "payload:", payload)

	// 변조
	parts := strings.Split(token, ".")
	tampered := parts[0] + "." + hex.EncodeToString([]byte("bad")) // 서명은 그대로 → 검증 실패
	_, ok = verify(key, tampered)
	fmt.Println("tampered ok:", ok)
}

func randomKey() []byte {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return b
}

func sign(key []byte, payload string) string {
	p := base64.RawURLEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(p))
	sig := mac.Sum(nil)

	return p + "." + hex.EncodeToString(sig)
}

func verify(key []byte, token string) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", false
	}
	p, sigHex := parts[0], parts[1]
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return "", false
	}

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(p))
	want := mac.Sum(nil)
	if !hmac.Equal(sig, want) {
		return "", false
	}

	b, err := base64.RawURLEncoding.DecodeString(p)
	if err != nil {
		return "", false
	}
	return string(b), true
}
