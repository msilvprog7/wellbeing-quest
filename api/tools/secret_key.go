package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	// Create random 32 bytes for 256 bits
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	// Encode to base64 for easy storage in .env
	encoded := base64.StdEncoding.EncodeToString(key)
	fmt.Println("Your 256-bit secret key (base64):", encoded)
}
