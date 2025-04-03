package main

import (
	"crypto/aes"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"os"
	"encoding/base64"
	"fmt"


	"github.com/joho/godotenv"

)
var err = godotenv.Load()
var encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))

func encryptID(text string) (string, error) {
	block, _ := aes.NewCipher(encryptionKey)
	nonce := make([]byte,12)
	io.ReadFull(rand.Reader, nonce)
	aesGCM, _ :=cipher.NewGCM(block)
	ciphertext := aesGCM.Seal(nil, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func decryptID(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	nonce := data[:12]         	
	ciphertext := data[12:]    
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plainText, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

