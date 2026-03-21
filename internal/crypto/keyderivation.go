package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize       = 32
	iterations     = 100000
	derivedKeySize = 32
)

func GenerateSalt() (string, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func DeriveKey(password, salt string) ([]byte, error) {
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if salt == "" {
		return nil, errors.New("salt cannot be empty")
	}

	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), saltBytes, iterations, derivedKeySize, sha256.New)
	return key, nil
}

func HashPassword(password string) (string, string, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return "", "", err
	}

	key, err := DeriveKey(password, salt)
	if err != nil {
		return "", "", err
	}

	hash := base64.StdEncoding.EncodeToString(key)
	return hash, salt, nil
}

func VerifyPassword(password, hash, salt string) bool {
	key, err := DeriveKey(password, salt)
	if err != nil {
		return false
	}

	expectedHash := base64.StdEncoding.EncodeToString(key)
	return hash == expectedHash
}
