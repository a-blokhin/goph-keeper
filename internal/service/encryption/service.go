package encryption

import (
	"github.com/a-blokhin/goph-keeper/internal/crypto"
	"github.com/a-blokhin/goph-keeper/internal/model"
)

type encryptionService struct {
	encryptor *crypto.Encryptor
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(encryptor *crypto.Encryptor) EncryptionService {
	return &encryptionService{
		encryptor: encryptor,
	}
}

// EncryptPassword encrypts a password string
func (s *encryptionService) EncryptPassword(password string) (string, error) {
	return s.encryptor.Encrypt([]byte(password))
}

// DecryptPassword decrypts an encrypted password
func (s *encryptionService) DecryptPassword(encryptedPassword string) (string, error) {
	decrypted, err := s.encryptor.Decrypt(encryptedPassword)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// EncryptText encrypts text data
func (s *encryptionService) EncryptText(text string) (string, error) {
	return s.encryptor.Encrypt([]byte(text))
}

// DecryptText decrypts encrypted text data
func (s *encryptionService) DecryptText(encryptedText string) (string, error) {
	decrypted, err := s.encryptor.Decrypt(encryptedText)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// EncryptBinary encrypts binary data
func (s *encryptionService) EncryptBinary(data []byte) ([]byte, error) {
	return s.encryptor.EncryptBytes(data)
}

// DecryptBinary decrypts encrypted binary data
func (s *encryptionService) DecryptBinary(encryptedData []byte) ([]byte, error) {
	return s.encryptor.DecryptBytes(encryptedData)
}

// EncryptCard encrypts card data
func (s *encryptionService) EncryptCard(cardNumber, cardHolder, expiry, cvv string) (encryptedCardNumber, encryptedCardHolder, encryptedExpiry, encryptedCVV string, err error) {
	encryptedCardNumber, err = s.encryptor.Encrypt([]byte(cardNumber))
	if err != nil {
		return "", "", "", "", err
	}

	encryptedCardHolder, err = s.encryptor.Encrypt([]byte(cardHolder))
	if err != nil {
		return "", "", "", "", err
	}

	encryptedExpiry, err = s.encryptor.Encrypt([]byte(expiry))
	if err != nil {
		return "", "", "", "", err
	}

	encryptedCVV, err = s.encryptor.Encrypt([]byte(cvv))
	if err != nil {
		return "", "", "", "", err
	}

	return encryptedCardNumber, encryptedCardHolder, encryptedExpiry, encryptedCVV, nil
}

// DecryptCard decrypts card data
func (s *encryptionService) DecryptCard(encryptedCardNumber, encryptedCardHolder, encryptedExpiry, encryptedCVV string) (cardNumber, cardHolder, expiry, cvv string, err error) {
	decryptedCardNumber, err := s.encryptor.Decrypt(encryptedCardNumber)
	if err != nil {
		return "", "", "", "", err
	}

	decryptedCardHolder, err := s.encryptor.Decrypt(encryptedCardHolder)
	if err != nil {
		return "", "", "", "", err
	}

	decryptedExpiry, err := s.encryptor.Decrypt(encryptedExpiry)
	if err != nil {
		return "", "", "", "", err
	}

	decryptedCVV, err := s.encryptor.Decrypt(encryptedCVV)
	if err != nil {
		return "", "", "", "", err
	}

	return string(decryptedCardNumber), string(decryptedCardHolder), string(decryptedExpiry), string(decryptedCVV), nil
}

// EncryptCredential encrypts credential password
func (s *encryptionService) EncryptCredential(credential *model.Credential) error {
	encrypted, err := s.encryptor.Encrypt([]byte(credential.PasswordEncrypted))
	if err != nil {
		return err
	}
	credential.PasswordEncrypted = encrypted
	return nil
}

// DecryptCredential decrypts credential password
func (s *encryptionService) DecryptCredential(credential *model.Credential) error {
	decrypted, err := s.encryptor.Decrypt(credential.PasswordEncrypted)
	if err != nil {
		return err
	}
	credential.PasswordEncrypted = string(decrypted)
	return nil
}

// EncryptTextData encrypts text data content
func (s *encryptionService) EncryptTextData(textData *model.TextData) error {
	encrypted, err := s.encryptor.Encrypt([]byte(textData.DataEncrypted))
	if err != nil {
		return err
	}
	textData.DataEncrypted = encrypted
	return nil
}

// DecryptTextData decrypts text data content
func (s *encryptionService) DecryptTextData(textData *model.TextData) error {
	decrypted, err := s.encryptor.Decrypt(textData.DataEncrypted)
	if err != nil {
		return err
	}
	textData.DataEncrypted = string(decrypted)
	return nil
}

// EncryptBinaryData encrypts binary data content
func (s *encryptionService) EncryptBinaryData(binaryData *model.BinaryData) error {
	encrypted, err := s.encryptor.EncryptBytes(binaryData.DataEncrypted)
	if err != nil {
		return err
	}
	binaryData.DataEncrypted = encrypted
	return nil
}

// DecryptBinaryData decrypts binary data content
func (s *encryptionService) DecryptBinaryData(binaryData *model.BinaryData) error {
	decrypted, err := s.encryptor.DecryptBytes(binaryData.DataEncrypted)
	if err != nil {
		return err
	}
	binaryData.DataEncrypted = decrypted
	return nil
}

// EncryptCardData encrypts card data fields
func (s *encryptionService) EncryptCardData(card *model.Card) error {
	encryptedCardNumber, err := s.encryptor.Encrypt([]byte(card.CardNumberEncrypted))
	if err != nil {
		return err
	}
	card.CardNumberEncrypted = encryptedCardNumber

	encryptedCardHolder, err := s.encryptor.Encrypt([]byte(card.CardHolderEncrypted))
	if err != nil {
		return err
	}
	card.CardHolderEncrypted = encryptedCardHolder

	encryptedExpiry, err := s.encryptor.Encrypt([]byte(card.ExpiryEncrypted))
	if err != nil {
		return err
	}
	card.ExpiryEncrypted = encryptedExpiry

	encryptedCVV, err := s.encryptor.Encrypt([]byte(card.CVVEncrypted))
	if err != nil {
		return err
	}
	card.CVVEncrypted = encryptedCVV

	return nil
}

// DecryptCardData decrypts card data fields
func (s *encryptionService) DecryptCardData(card *model.Card) error {
	decryptedCardNumber, err := s.encryptor.Decrypt(card.CardNumberEncrypted)
	if err != nil {
		return err
	}
	card.CardNumberEncrypted = string(decryptedCardNumber)

	decryptedCardHolder, err := s.encryptor.Decrypt(card.CardHolderEncrypted)
	if err != nil {
		return err
	}
	card.CardHolderEncrypted = string(decryptedCardHolder)

	decryptedExpiry, err := s.encryptor.Decrypt(card.ExpiryEncrypted)
	if err != nil {
		return err
	}
	card.ExpiryEncrypted = string(decryptedExpiry)

	decryptedCVV, err := s.encryptor.Decrypt(card.CVVEncrypted)
	if err != nil {
		return err
	}
	card.CVVEncrypted = string(decryptedCVV)

	return nil
}
