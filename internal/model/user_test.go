package model

import (
	"testing"
	"time"
)

func TestUser_Fields(t *testing.T) {
	user := &User{
		ID:           "user123",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Salt:         "salt123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if user.ID != "user123" {
		t.Errorf("User.ID = %v, want user123", user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("User.Email = %v, want test@example.com", user.Email)
	}

	if user.PasswordHash != "hashedpassword" {
		t.Errorf("User.PasswordHash = %v, want hashedpassword", user.PasswordHash)
	}

	if user.Salt != "salt123" {
		t.Errorf("User.Salt = %v, want salt123", user.Salt)
	}
}

func TestUser_EmptyFields(t *testing.T) {
	user := &User{}

	if user.ID != "" {
		t.Errorf("User.ID should be empty, got %v", user.ID)
	}

	if user.Email != "" {
		t.Errorf("User.Email should be empty, got %v", user.Email)
	}

	if user.PasswordHash != "" {
		t.Errorf("User.PasswordHash should be empty, got %v", user.PasswordHash)
	}

	if user.Salt != "" {
		t.Errorf("User.Salt should be empty, got %v", user.Salt)
	}
}

func TestUser_Timestamps(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:           "user123",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Salt:         "salt123",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if user.CreatedAt.IsZero() {
		t.Error("User.CreatedAt is zero")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("User.UpdatedAt is zero")
	}

	if user.CreatedAt.After(now.Add(time.Second)) {
		t.Error("User.CreatedAt is in the future")
	}

	if user.UpdatedAt.After(now.Add(time.Second)) {
		t.Error("User.UpdatedAt is in the future")
	}
}

func TestUser_EmailValidation(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			valid: true,
		},
		{
			name:  "valid email with subdomain",
			email: "test@mail.example.com",
			valid: true,
		},
		{
			name:  "empty email",
			email: "",
			valid: false,
		},
		{
			name:  "email without @",
			email: "testexample.com",
			valid: false,
		},
		{
			name:  "email without domain",
			email: "test@",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Email: tt.email,
			}

			atIndex := -1
			for i, char := range user.Email {
				if char == '@' {
					atIndex = i
					break
				}
			}

			hasAt := atIndex != -1
			notAtStart := atIndex > 0
			notAtEnd := atIndex < len(user.Email)-1

			isValid := user.Email != "" && len(user.Email) > 3 && hasAt && notAtStart && notAtEnd
			if isValid != tt.valid {
				t.Errorf("User.Email validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}
