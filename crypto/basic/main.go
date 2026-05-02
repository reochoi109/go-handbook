package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	randomDemo()
	hashDemo()
	hmacDemo()
}

func randomDemo() {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	fmt.Println("err:", err)
	fmt.Println("random hex:", hex.EncodeToString(b))
}

func hashDemo() {
	sum := sha256.Sum256([]byte("hello"))
	fmt.Printf("sha256(hello)=%x\n", sum[:])
}

func hmacDemo() {
	key := []byte("secret-key")
	msg := []byte("payload")

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(msg)
	tag := mac.Sum(nil)
	fmt.Printf("hmac=%x\n", tag)

	// verify
	mac2 := hmac.New(sha256.New, key)
	_, _ = mac2.Write(msg)
	tag2 := mac2.Sum(nil)

	fmt.Println("valid:", hmac.Equal(tag, tag2))
}
