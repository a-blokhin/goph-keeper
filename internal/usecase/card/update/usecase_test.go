package update

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

func TestUpdateCardUsecase_Execute_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	existingCard := &model.Card{
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

	updatedCard := &model.Card{
		ID:      "card-123",
		Title:   "Updated Card",
		Meta:    "updated personal card",
		Version: 1,
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	// Mock encryption service to return success
	encService.EXPECT().EncryptCardData(updatedCard).Return(nil)

	// Mock repository to return success on update
	cardRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(c *model.Card) bool {
		return c.ID == "card-123" &&
			c.UserID == "user-123" &&
			c.Title == "Updated Card" &&
			c.Meta == "updated personal card" &&
			c.Version == 2 // Version should be incremented
	})).Return(nil)

	err := uc.Execute(context.Background(), "user-123", updatedCard)

	assert.NoError(t, err)
}

func TestUpdateCardUsecase_Execute_CardNotFound(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	card := &model.Card{
		ID:      "card-123",
		Title:   "Updated Card",
		Version: 1,
	}

	// Mock repository to return an error
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(nil, model.ErrCardNotFound)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
	assert.Equal(t, model.ErrCardNotFound, err)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "EncryptCardData")
	// Repository update should not be called
	cardRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCardUsecase_Execute_UnauthorizedAccess(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	existingCard := &model.Card{
		ID:      "card-123",
		UserID:  "user-456", // Different user
		Title:   "My Card",
		Version: 1,
	}

	card := &model.Card{
		ID:      "card-123",
		Title:   "Updated Card",
		Version: 1,
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
	assert.Equal(t, model.ErrForbidden, err)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "EncryptCardData")
	// Repository update should not be called
	cardRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCardUsecase_Execute_VersionConflict(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	existingCard := &model.Card{
		ID:      "card-123",
		UserID:  "user-123",
		Title:   "My Card",
		Version: 2, // Different version
	}

	card := &model.Card{
		ID:      "card-123",
		Title:   "Updated Card",
		Version: 1, // Conflicting version
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
	assert.Equal(t, model.ErrVersionConflict, err)
	// Encryption service should not be called
	encService.AssertNotCalled(t, "EncryptCardData")
	// Repository update should not be called
	cardRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCardUsecase_Execute_EncryptionError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	existingCard := &model.Card{
		ID:      "card-123",
		UserID:  "user-123",
		Title:   "My Card",
		Version: 1,
	}

	card := &model.Card{
		ID:      "card-123",
		Title:   "Updated Card",
		Version: 1,
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	// Mock encryption service to return an error
	encService.EXPECT().EncryptCardData(card).Return(assert.AnError)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
	// Repository update should not be called
	cardRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCardUsecase_Execute_RepositoryError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	existingCard := &model.Card{
		ID:      "card-123",
		UserID:  "user-123",
		Title:   "My Card",
		Version: 1,
	}

	card := &model.Card{
		ID:      "card-123",
		Title:   "Updated Card",
		Version: 1,
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	// Mock encryption service to return success
	encService.EXPECT().EncryptCardData(card).Return(nil)

	// Mock repository to return an error on update
	cardRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(assert.AnError)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
}
