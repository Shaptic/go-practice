package main

// Trying out the crypto package!

import (
	"os"
	"fmt"
	"flag"
	"strings"
	// "strconv"

	"crypto/aes"
	"crypto/sha1"
	"crypto/ed25519"
)

func isUpperLetter(c int) bool {
	return (c >= 65 && c <= 90)
}

func isLowerLetter(c int) bool {
	return (c >= 97 && c <= 122)
}

func rot13(plaintext string) string {
	return caesar(plaintext, 13)
}

func caesar(plaintext string, amount int) string {
	var ciphertext strings.Builder

	for _, char := range plaintext {
		c := int(char)
		if isUpperLetter(c) {
			c = ((c - 65 + amount) % 26) + 65
		} else if isLowerLetter(c) {
			c = ((c - 97 + amount) % 26) + 97
		}
		ciphertext.WriteString(string(c))
	}

	return ciphertext.String()
}

func proper(plaintext string) (ciphertext string) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Printf("Failed to generate Ed25519 keypair: %s", err)
		return
	}

	hasher := sha1.New()
	key := hasher.Sum([]byte("this is a super secure key lol"))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./ciphers [plaintext to encipher]")
		return
	}

	key := flag.Int("amount", 13, "amount to rotate text")
	algo := flag.String("algo", "caesar", "type of 'encryption'")
	flag.Parse()
	plaintext := strings.Join(flag.Args(), " ")

	var ciphertext string
	if *algo == "caesar" {
		ciphertext = caesar(plaintext, *key)
	}

	fmt.Println(ciphertext)
}
