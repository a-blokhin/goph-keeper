package list

import (
	"context"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestListCardsUsecase_Execute_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	cards := []*model.Card{
		{
			ID:                  "card-123",
			UserID:              "user-123",
			Title:               "My Card",
			CardNumberEncrypted: "encrypted-card-number-1",
			CardHolderEncrypted: "encrypted-card-holder-1",
			ExpiryEncrypted:     "encrypted-expiry-1",
			CVVEncrypted:        "encrypted-cvv-1",
			Meta:                "personal card",
			Version:             1,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  "card-456",
			UserID:              "user-123",
			Title:               "Another Card",
			CardNumberEncrypted: "encrypted-card-number-2",
			CardHolderEncrypted: "encrypted-card-holder-2",
			ExpiryEncrypted:     "encrypted-expiry-2",
			CVVEncrypted:        "encrypted-cvv-2",
			Meta:                "business card",
			Version:             1,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	// Mock repository to return the cards
	cardRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(cards, nil)

	// Mock encryption service to return success for each card
	encService.EXPECT().DecryptCardData(cards[0]).Return(nil)
	encService.EXPECT().DecryptCardData(cards[1]).Return(nil)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Equal(t, cards, result)
}

func TestListCardsUsecase_Execute_RepositoryError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	// Mock repository to return an error
	cardRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(nil, assert.AnError)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "DecryptCardData")
}

func TestListCardsUsecase_Execute_DecryptionError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	cards := []*model.Card{
		{
			ID:                  "card-123",
			UserID:              "user-123",
			Title:               "My Card",
			CardNumberEncrypted: "encrypted-card-number-1",
			CardHolderEncrypted: "encrypted-card-holder-1",
			ExpiryEncrypted:     "encrypted-expiry-1",
			CVVEncrypted:        "encrypted-cvv-1",
			Meta:                "personal card",
			Version:             1,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  "card-456",
			UserID:              "user-123",
			Title:               "Another Card",
			CardNumberEncrypted: "encrypted-card-number-2",
			CardHolderEncrypted: "encrypted-card-holder-2",
			ExpiryEncrypted:     "encrypted-expiry-2",
			CVVEncrypted:        "encrypted-cvv-2",
			Meta:                "business card",
			Version:             1,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	// Mock repository to return the cards
	cardRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(cards, nil)

	// Mock encryption service to return success for the first card but error for the second
	encService.EXPECT().DecryptCardData(cards[0]).Return(nil)
	encService.EXPECT().DecryptCardData(cards[1]).Return(assert.AnError)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListCardsUsecase_Execute_EmptyList(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	// Mock repository to return an empty list
	cardRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return([]*model.Card{}, nil)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Empty(t, result)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "DecryptCardData")
}
