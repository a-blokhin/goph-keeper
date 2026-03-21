package model

import (
	"testing"
	"time"
)

func TestCredential_Fields(t *testing.T) {
	cred := &Credential{
		ID:                "cred123",
		UserID:            "user123",
		Title:             "Test Credential",
		Login:             "testuser",
		PasswordEncrypted: "encryptedpassword",
		Meta:              "test meta",
		Version:           1,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if cred.ID != "cred123" {
		t.Errorf("Credential.ID = %v, want cred123", cred.ID)
	}

	if cred.UserID != "user123" {
		t.Errorf("Credential.UserID = %v, want user123", cred.UserID)
	}

	if cred.Title != "Test Credential" {
		t.Errorf("Credential.Title = %v, want Test Credential", cred.Title)
	}

	if cred.Login != "testuser" {
		t.Errorf("Credential.Login = %v, want testuser", cred.Login)
	}

	if cred.PasswordEncrypted != "encryptedpassword" {
		t.Errorf("Credential.PasswordEncrypted = %v, want encryptedpassword", cred.PasswordEncrypted)
	}

	if cred.Meta != "test meta" {
		t.Errorf("Credential.Meta = %v, want test meta", cred.Meta)
	}

	if cred.Version != 1 {
		t.Errorf("Credential.Version = %v, want 1", cred.Version)
	}
}

func TestCredential_VersionConflict(t *testing.T) {
	tests := []struct {
		name         string
		clientVer    int32
		serverVer    int32
		wantConflict bool
	}{
		{
			name:         "no conflict - same version",
			clientVer:    1,
			serverVer:    1,
			wantConflict: false,
		},
		{
			name:         "conflict - client behind",
			clientVer:    1,
			serverVer:    2,
			wantConflict: true,
		},
		{
			name:         "conflict - client ahead",
			clientVer:    3,
			serverVer:    2,
			wantConflict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := &Credential{
				Version: tt.serverVer,
			}
			got := cred.Version != tt.clientVer
			if got != tt.wantConflict {
				t.Errorf("Credential version conflict = %v, want %v", got, tt.wantConflict)
			}
		})
	}
}

func TestCredential_Timestamps(t *testing.T) {
	now := time.Now()
	cred := &Credential{
		ID:                "cred123",
		UserID:            "user123",
		Title:             "Test Credential",
		Login:             "testuser",
		PasswordEncrypted: "encryptedpassword",
		Version:           1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if cred.CreatedAt.IsZero() {
		t.Error("Credential.CreatedAt is zero")
	}

	if cred.UpdatedAt.IsZero() {
		t.Error("Credential.UpdatedAt is zero")
	}

	if cred.CreatedAt.After(now.Add(time.Second)) {
		t.Error("Credential.CreatedAt is in the future")
	}

	if cred.UpdatedAt.After(now.Add(time.Second)) {
		t.Error("Credential.UpdatedAt is in the future")
	}
}

func TestCredential_EmptyFields(t *testing.T) {
	cred := &Credential{}

	if cred.ID != "" {
		t.Errorf("Credential.ID should be empty, got %v", cred.ID)
	}

	if cred.UserID != "" {
		t.Errorf("Credential.UserID should be empty, got %v", cred.UserID)
	}

	if cred.Title != "" {
		t.Errorf("Credential.Title should be empty, got %v", cred.Title)
	}

	if cred.Login != "" {
		t.Errorf("Credential.Login should be empty, got %v", cred.Login)
	}

	if cred.PasswordEncrypted != "" {
		t.Errorf("Credential.PasswordEncrypted should be empty, got %v", cred.PasswordEncrypted)
	}

	if cred.Version != 0 {
		t.Errorf("Credential.Version should be 0, got %v", cred.Version)
	}
}
