package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"strconv"
)

// NewKey returns a new key for Encrypt/Decrypt
func NewKey(size int) ([]byte, error) {

	key := make([]byte, size)

	switch size {
	default:
		return key, errors.New("invalid key size " + strconv.Itoa(size))
	case 16, 24, 32:
		break
	}

	// create random key
	_, err := rand.Read(key)
	if err != nil {
		return key, err
	}

	return key, nil
}

// Encrypt returns an encrypted byte array
func Encrypt(key []byte, text []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

// Decrypt returns a decrypted byte array
func Decrypt(key []byte, ciphertext []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	if len(ciphertext) == 0 {
		return []byte{}, errors.New("no text")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return []byte{}, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}

	return plaintext, nil
}
