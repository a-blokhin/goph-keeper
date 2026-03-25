package model

import "time"

type Credential struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Title             string    `json:"title"`
	Login             string    `json:"login"`
	PasswordEncrypted string    `json:"password_encrypted"`
	Meta              string    `json:"meta"`
	Version           int32     `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
