package crypto

import (
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error = %v", err)
	}

	if len(salt1) == 0 {
		t.Error("GenerateSalt() returned empty salt")
	}

	salt2, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error = %v", err)
	}

	if salt1 == salt2 {
		t.Error("GenerateSalt() returned same salt twice, salts should be unique")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, salt, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	if salt == "" {
		t.Error("HashPassword() returned empty salt")
	}

	hash2, salt2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == hash2 {
		t.Error("HashPassword() returned same hash for same password, hashes should be different due to salt")
	}

	if salt == salt2 {
		t.Error("HashPassword() returned same salt twice, salts should be unique")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"

	hash, salt, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	tests := []struct {
		name string
		hash string
		salt string
		want bool
	}{
		{
			name: "correct password",
			hash: hash,
			salt: salt,
			want: true,
		},
		{
			name: "incorrect password",
			hash: hash,
			salt: salt,
			want: false,
		},
		{
			name: "empty password",
			hash: hash,
			salt: salt,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testPassword string
			if tt.name == "correct password" {
				testPassword = password
			} else if tt.name == "incorrect password" {
				testPassword = "wrongpassword"
			} else {
				testPassword = ""
			}

			got := VerifyPassword(testPassword, tt.hash, tt.salt)
			if got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeriveKey(t *testing.T) {
	password := "testpassword123"
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error = %v", err)
	}

	key1, err := DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	if len(key1) != derivedKeySize {
		t.Errorf("DeriveKey() returned key of length %d, want %d", len(key1), derivedKeySize)
	}

	key2, err := DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	if string(key1) != string(key2) {
		t.Error("DeriveKey() returned different keys for same password and salt")
	}

	salt2, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error = %v", err)
	}

	key3, err := DeriveKey(password, salt2)
	if err != nil {
		t.Fatalf("DeriveKey() error = %v", err)
	}

	if string(key1) == string(key3) {
		t.Error("DeriveKey() returned same key for different salts")
	}
}

func TestHashPasswordVerifyPasswordRoundtrip(t *testing.T) {
	passwords := []string{
		"simple",
		"complex!@#$%^&*()",
		"with spaces",
		"with numbers 12345",
		"very long password with many characters and symbols !@#$%^&*()_+-=[]{}|;':,.<>?",
	}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			hash, salt, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			valid := VerifyPassword(password, hash, salt)
			if !valid {
				t.Error("VerifyPassword() returned false for correct password")
			}

			invalidValid := VerifyPassword("wrong"+password, hash, salt)
			if invalidValid {
				t.Error("VerifyPassword() returned true for incorrect password")
			}
		})
	}
}
