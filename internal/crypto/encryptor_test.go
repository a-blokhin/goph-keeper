package crypto

import (
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid key",
			key:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "invalid key too short",
			key:     make([]byte, 16),
			wantErr: true,
		},
		{
			name:    "invalid key too long",
			key:     make([]byte, 64),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEncryptor(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptor_Encrypt(t *testing.T) {
	key := make([]byte, 32)
	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "encrypt empty data",
			data:    []byte{},
			wantErr: false,
		},
		{
			name:    "encrypt small data",
			data:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "encrypt large data",
			data:    make([]byte, 1024),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encryptor.Encrypt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Error("Encrypt() returned empty string")
			}
		})
	}
}

func TestEncryptor_Decrypt(t *testing.T) {
	key := make([]byte, 32)
	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	originalData := []byte("test data")
	encrypted, err := encryptor.Encrypt(originalData)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	tests := []struct {
		name    string
		data    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "decrypt valid data",
			data:    encrypted,
			want:    originalData,
			wantErr: false,
		},
		{
			name:    "decrypt invalid base64",
			data:    "invalid!!!",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "decrypt too short data",
			data:    "YWJj",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encryptor.Decrypt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(result) != string(tt.want) {
				t.Errorf("Decrypt() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestEncryptor_EncryptBytes(t *testing.T) {
	key := make([]byte, 32)
	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "encrypt empty bytes",
			data:    []byte{},
			wantErr: false,
		},
		{
			name:    "encrypt small bytes",
			data:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "encrypt large bytes",
			data:    make([]byte, 1024),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encryptor.EncryptBytes(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result) == 0 {
				t.Error("EncryptBytes() returned empty bytes")
			}
		})
	}
}

func TestEncryptor_EncryptDecryptRoundtrip(t *testing.T) {
	key := make([]byte, 32)
	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	testData := []byte("This is a test message for encryption and decryption")

	encrypted, err := encryptor.Encrypt(testData)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data doesn't match original: got %v, want %v", decrypted, testData)
	}
}

func TestEncryptor_EncryptBytesDecryptRoundtrip(t *testing.T) {
	key := make([]byte, 32)
	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	testData := []byte("This is a test message for bytes encryption and decryption")

	encrypted, err := encryptor.EncryptBytes(testData)
	if err != nil {
		t.Fatalf("Failed to encrypt bytes: %v", err)
	}

	decrypted, err := encryptor.DecryptBytes(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt bytes: %v", err)
	}

	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted bytes don't match original: got %v, want %v", decrypted, testData)
	}
}
