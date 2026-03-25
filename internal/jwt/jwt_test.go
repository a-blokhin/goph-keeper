package jwt

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	secretKey := "test-secret-key-12345678901234567890"
	manager := NewJWTManager(secretKey)

	userID := "user123"
	email := "test@example.com"

	token, err := manager.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	secretKey := "test-secret-key-12345678901234567890"
	manager := NewJWTManager(secretKey)

	userID := "user123"
	email := "test@example.com"

	token, err := manager.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateToken() UserID = %v, want %v", claims.UserID, userID)
	}

	if claims.Email != email {
		t.Errorf("ValidateToken() Email = %v, want %v", claims.Email, email)
	}
}

func TestJWTManager_ValidateToken_Invalid(t *testing.T) {
	secretKey := "test-secret-key-12345678901234567890"
	manager := NewJWTManager(secretKey)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not.a.valid.jwt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	secretKey1 := "test-secret-key-12345678901234567890"
	manager1 := NewJWTManager(secretKey1)

	secretKey2 := "different-secret-key-123456789012345"
	manager2 := NewJWTManager(secretKey2)

	userID := "user123"
	email := "test@example.com"

	token, err := manager1.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken() expected error for token with wrong secret")
	}
}

func TestJWTManager_TokenExpiration(t *testing.T) {
	secretKey := "test-secret-key-12345678901234567890"
	manager := NewJWTManager(secretKey)

	userID := "user123"
	email := "test@example.com"

	token, err := manager.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	expectedExpiry := time.Now().Add(24 * time.Hour)
	timeDiff := expectedExpiry.Sub(claims.ExpiresAt.Time)

	if timeDiff < -time.Minute || timeDiff > time.Minute {
		t.Errorf("ValidateToken() ExpiresAt = %v, want approximately %v", claims.ExpiresAt.Time, expectedExpiry)
	}
}
