package encription

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetKey() string {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	key := os.Getenv("DB_HOST")

	return key
}

// Encrypt encrypts plaintext using AES-GCM.
func EncryptData(plaintext, key string) (string, error) {

	// Convert the key to bytes and validate its length
	keyBytes := []byte(key)
	if len(keyBytes) != 32 { // AES-256 requires a 32-byte key
		return "", errors.New("invalid key length: must be 32 bytes")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a nonce of the correct size
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to hex for easy storage or transfer
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM.
func DecryptData(ciphertext, key string) (string, error) {

	// Convert the key and ciphertext to bytes
	keyBytes := []byte(key)
	if len(keyBytes) != 32 { // AES-256 requires a 32-byte key
		return "", errors.New("invalid key length: must be 32 bytes")
	}

	cipherBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Split the nonce and the actual ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
