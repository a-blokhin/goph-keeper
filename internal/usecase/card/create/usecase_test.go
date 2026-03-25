package create

import (
	"context"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateCardUsecase_Execute_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	card := &model.Card{
		Title: "My Card",
		Meta:  "personal card",
	}

	// Mock encryption service to return success
	encService.EXPECT().EncryptCardData(card).Return(nil)

	// Mock repository to return success
	cardRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(c *model.Card) bool {
		return c.UserID == "user-123" &&
			c.Title == "My Card" &&
			c.Meta == "personal card" &&
			c.Version == 1
	})).Return(nil)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.NoError(t, err)
}

func TestCreateCardUsecase_Execute_EncryptionError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	card := &model.Card{
		Title: "My Card",
		Meta:  "personal card",
	}

	// Mock encryption service to return an error
	encService.EXPECT().EncryptCardData(card).Return(assert.AnError)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
	// Repository should not be called
	cardRepo.AssertNotCalled(t, "Create")
}

func TestCreateCardUsecase_Execute_RepositoryError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()
	uc := New(cardRepo, encService, logger)

	card := &model.Card{
		Title: "My Card",
		Meta:  "personal card",
	}

	// Mock encryption service to return success
	encService.EXPECT().EncryptCardData(card).Return(nil)

	// Mock repository to return an error
	cardRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(assert.AnError)

	err := uc.Execute(context.Background(), "user-123", card)

	assert.Error(t, err)
}
