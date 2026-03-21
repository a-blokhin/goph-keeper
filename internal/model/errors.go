package model

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidToken       = errors.New("invalid token")

	ErrCredentialNotFound      = errors.New("credential not found")
	ErrCredentialAlreadyExists = errors.New("credential already exists")
	ErrVersionConflict         = errors.New("version conflict: data has been modified by another client")

	ErrTextDataNotFound      = errors.New("text data not found")
	ErrTextDataAlreadyExists = errors.New("text data already exists")

	ErrBinaryDataNotFound      = errors.New("binary data not found")
	ErrBinaryDataAlreadyExists = errors.New("binary data already exists")
	ErrBinaryDataTooLarge      = errors.New("binary data exceeds 10MB limit")

	ErrCardNotFound      = errors.New("card not found")
	ErrCardAlreadyExists = errors.New("card already exists")

	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)
