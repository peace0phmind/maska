package mask

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"math/big"
)

// Base62 encoding characters
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// base62Encode encodes a byte slice to a Base62 string
func base62Encode(data []byte) string {
	var result string
	num := new(big.Int).SetBytes(data)

	// Convert to Base62
	for num.Cmp(big.NewInt(0)) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, big.NewInt(62), remainder)
		result = string(base62Chars[remainder.Int64()]) + result
	}

	return result
}

// Pad applies ANSI X.923 padding to the plaintext
func pad(plaintext []byte) []byte {
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := make([]byte, len(plaintext)+padding)
	copy(padtext, plaintext)
	for i := len(plaintext); i < len(padtext)-1; i++ {
		padtext[i] = 0
	}
	padtext[len(padtext)-1] = byte(padding)
	return padtext
}

// AESEncrypt encrypts a byte slice using AES-128 in CBC mode with ANSI X.923 padding
func AESEncrypt(key []byte, plaintext []byte) (string, error) {
	if len(key) != 16 {
		return "", errors.New("key length must be 16 bytes for AES-128")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Use the new padding function
	padtext := pad(plaintext)

	// Generate random IV
	ciphertext := make([]byte, aes.BlockSize+len(padtext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padtext)

	return base62Encode(ciphertext), nil
}

// AESDecrypt decrypts a string that was encrypted using AESEncrypt
func AESDecrypt(key []byte, encrypted string) ([]byte, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Extract IV
	if len(ciphertext) < aes.BlockSize {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	padding := int(ciphertext[len(ciphertext)-1])
	if padding > len(ciphertext) {
		return nil, err
	}

	return ciphertext[:len(ciphertext)-padding], nil
}
