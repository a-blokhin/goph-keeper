package cardget

import (
	"context"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCardUsecase_Execute_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	uc := New(cardRepo, encService)

	expectedCard := &model.Card{
		ID:                  "card-123",
		UserID:              "user-123",
		Title:               "My Card",
		CardNumberEncrypted: "encrypted-card-number",
		CardHolderEncrypted: "encrypted-card-holder",
		ExpiryEncrypted:     "encrypted-expiry",
		CVVEncrypted:        "encrypted-cvv",
		Meta:                "personal card",
		Version:             1,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Mock repository to return the card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(expectedCard, nil)

	// Mock encryption service to return success
	encService.EXPECT().DecryptCardData(expectedCard).Return(nil)

	card, err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestGetCardUsecase_Execute_CardNotFound(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	uc := New(cardRepo, encService)

	// Mock repository to return an error
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(nil, model.ErrCardNotFound)

	card, err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.Error(t, err)
	assert.Equal(t, model.ErrCardNotFound, err)
	assert.Nil(t, card)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "DecryptCardData")
}

func TestGetCardUsecase_Execute_UnauthorizedAccess(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	uc := New(cardRepo, encService)

	card := &model.Card{
		ID:     "card-123",
		UserID: "user-456", // Different user
		Title:  "My Card",
	}

	// Mock repository to return the card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(card, nil)

	card, err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.Error(t, err)
	assert.Equal(t, model.ErrCardNotFound, err)
	assert.Nil(t, card)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "DecryptCardData")
}

func TestGetCardUsecase_Execute_DecryptionError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	uc := New(cardRepo, encService)

	card := &model.Card{
		ID:     "card-123",
		UserID: "user-123",
		Title:  "My Card",
	}

	// Mock repository to return the card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(card, nil)

	// Mock encryption service to return an error
	encService.EXPECT().DecryptCardData(card).Return(assert.AnError)

	card, err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.Error(t, err)
	assert.Nil(t, card)
}
