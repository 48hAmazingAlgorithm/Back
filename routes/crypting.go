package routes

import (
	"io"
	"os"

	"crypto/aes"
	"crypto/rand"
	"crypto/cipher"
	
	"encoding/base64"
	
	"github.com/joho/godotenv"
)

var err = godotenv.Load()
var encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY"))

func EncryptID(text string) (string, error) {
	block, _ := aes.NewCipher(encryptionKey)
	nonce := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)
	aesGCM, _ := cipher.NewGCM(block)
	ciphertext := aesGCM.Seal(nil, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func DecryptID(text string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(text)
	nonce := data[:12]
	ciphertext := data[12:]
	block, _ := aes.NewCipher(encryptionKey)
	aesGCM, _ := cipher.NewGCM(block)
	plainText, _ := aesGCM.Open(nil, nonce, ciphertext, nil)
	return string(plainText), nil
}