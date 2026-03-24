package encryption

import (
	"github.com/a-blokhin/goph-keeper/internal/model"
)

// EncryptionService provides encryption and decryption operations for sensitive data
type EncryptionService interface {
	// EncryptPassword encrypts a password string
	EncryptPassword(password string) (string, error)

	// DecryptPassword decrypts an encrypted password
	DecryptPassword(encryptedPassword string) (string, error)

	// EncryptText encrypts text data
	EncryptText(text string) (string, error)

	// DecryptText decrypts encrypted text data
	DecryptText(encryptedText string) (string, error)

	// EncryptBinary encrypts binary data
	EncryptBinary(data []byte) ([]byte, error)

	// DecryptBinary decrypts encrypted binary data
	DecryptBinary(encryptedData []byte) ([]byte, error)

	// EncryptCard encrypts card data
	EncryptCard(cardNumber, cardHolder, expiry, cvv string) (encryptedCardNumber, encryptedCardHolder, encryptedExpiry, encryptedCVV string, err error)

	// DecryptCard decrypts card data
	DecryptCard(encryptedCardNumber, encryptedCardHolder, encryptedExpiry, encryptedCVV string) (cardNumber, cardHolder, expiry, cvv string, err error)

	// EncryptCredential encrypts credential password
	EncryptCredential(credential *model.Credential) error

	// DecryptCredential decrypts credential password
	DecryptCredential(credential *model.Credential) error

	// EncryptTextData encrypts text data content
	EncryptTextData(textData *model.TextData) error

	// DecryptTextData decrypts text data content
	DecryptTextData(textData *model.TextData) error

	// EncryptBinaryData encrypts binary data content
	EncryptBinaryData(binaryData *model.BinaryData) error

	// DecryptBinaryData decrypts binary data content
	DecryptBinaryData(binaryData *model.BinaryData) error

	// EncryptCardData encrypts card data fields
	EncryptCardData(card *model.Card) error

	// DecryptCardData decrypts card data fields
	DecryptCardData(card *model.Card) error
}
