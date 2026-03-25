package encryption

import (
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/crypto"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptionService_EncryptPassword(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	encrypted, err := service.EncryptPassword("mysecretpassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	// Verify we can decrypt it back
	decrypted, err := service.DecryptPassword(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, "mysecretpassword", decrypted)
}

func TestEncryptionService_DecryptPassword_InvalidData(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	// Try to decrypt invalid data
	_, err = service.DecryptPassword("invalid-base64")
	assert.Error(t, err)
}

func TestEncryptionService_EncryptText(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	encrypted, err := service.EncryptText("my secret text")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	// Verify we can decrypt it back
	decrypted, err := service.DecryptText(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, "my secret text", decrypted)
}

func TestEncryptionService_EncryptBinary(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	data := []byte("my binary data")
	encrypted, err := service.EncryptBinary(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	// Verify we can decrypt it back
	decrypted, err := service.DecryptBinary(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted)
}

func TestEncryptionService_EncryptCard(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	encryptedCardNumber, encryptedCardHolder, encryptedExpiry, encryptedCVV, err := service.EncryptCard(
		"1234567890123456",
		"John Doe",
		"12/25",
		"123",
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedCardNumber)
	assert.NotEmpty(t, encryptedCardHolder)
	assert.NotEmpty(t, encryptedExpiry)
	assert.NotEmpty(t, encryptedCVV)

	// Verify we can decrypt it back
	cardNumber, cardHolder, expiry, cvv, err := service.DecryptCard(
		encryptedCardNumber,
		encryptedCardHolder,
		encryptedExpiry,
		encryptedCVV,
	)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890123456", cardNumber)
	assert.Equal(t, "John Doe", cardHolder)
	assert.Equal(t, "12/25", expiry)
	assert.Equal(t, "123", cvv)
}

func TestEncryptionService_EncryptCredential(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	credential := &model.Credential{
		PasswordEncrypted: "mysecretpassword",
	}

	err = service.EncryptCredential(credential)
	assert.NoError(t, err)
	assert.NotEmpty(t, credential.PasswordEncrypted)
	assert.NotEqual(t, "mysecretpassword", credential.PasswordEncrypted)

	// Verify we can decrypt it back
	err = service.DecryptCredential(credential)
	assert.NoError(t, err)
	assert.Equal(t, "mysecretpassword", credential.PasswordEncrypted)
}

func TestEncryptionService_EncryptTextData(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	textData := &model.TextData{
		DataEncrypted: "my secret text data",
	}

	err = service.EncryptTextData(textData)
	assert.NoError(t, err)
	assert.NotEmpty(t, textData.DataEncrypted)
	assert.NotEqual(t, "my secret text data", textData.DataEncrypted)

	// Verify we can decrypt it back
	err = service.DecryptTextData(textData)
	assert.NoError(t, err)
	assert.Equal(t, "my secret text data", textData.DataEncrypted)
}

func TestEncryptionService_EncryptBinaryData(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	binaryData := &model.BinaryData{
		DataEncrypted: []byte("my binary data"),
	}

	err = service.EncryptBinaryData(binaryData)
	assert.NoError(t, err)
	assert.NotEmpty(t, binaryData.DataEncrypted)
	assert.NotEqual(t, []byte("my binary data"), binaryData.DataEncrypted)

	// Verify we can decrypt it back
	err = service.DecryptBinaryData(binaryData)
	assert.NoError(t, err)
	assert.Equal(t, []byte("my binary data"), binaryData.DataEncrypted)
}

func TestEncryptionService_EncryptCardData(t *testing.T) {
	// Create a valid 32-byte key for testing
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)

	service := NewEncryptionService(encryptor)

	card := &model.Card{
		CardNumberEncrypted: "1234567890123456",
		CardHolderEncrypted: "John Doe",
		ExpiryEncrypted:     "12/25",
		CVVEncrypted:        "123",
	}

	err = service.EncryptCardData(card)
	assert.NoError(t, err)
	assert.NotEmpty(t, card.CardNumberEncrypted)
	assert.NotEmpty(t, card.CardHolderEncrypted)
	assert.NotEmpty(t, card.ExpiryEncrypted)
	assert.NotEmpty(t, card.CVVEncrypted)
	assert.NotEqual(t, "1234567890123456", card.CardNumberEncrypted)
	assert.NotEqual(t, "John Doe", card.CardHolderEncrypted)
	assert.NotEqual(t, "12/25", card.ExpiryEncrypted)
	assert.NotEqual(t, "123", card.CVVEncrypted)

	// Verify we can decrypt it back
	err = service.DecryptCardData(card)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890123456", card.CardNumberEncrypted)
	assert.Equal(t, "John Doe", card.CardHolderEncrypted)
	assert.Equal(t, "12/25", card.ExpiryEncrypted)
	assert.Equal(t, "123", card.CVVEncrypted)
}
