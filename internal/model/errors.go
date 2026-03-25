package model

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrCredentialNotFound = errors.New("credential not found")

	ErrTextDataNotFound = errors.New("text data not found")

	ErrBinaryDataNotFound = errors.New("binary data not found")
	ErrBinaryDataTooLarge = errors.New("binary data exceeds 10MB limit")

	ErrCardNotFound = errors.New("card not found")

	ErrVersionConflict = errors.New("version conflict: data has been modified by another client")

	ErrForbidden = errors.New("forbidden")
)
