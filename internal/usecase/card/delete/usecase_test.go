package delete

import (
	"context"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestDeleteCardUsecase_Execute_Success(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	logger := zap.NewNop()
	uc := New(cardRepo, logger)

	existingCard := &model.Card{
		ID:        "card-123",
		UserID:    "user-123",
		Title:     "My Card",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	// Mock repository to return success on delete
	cardRepo.EXPECT().Delete(mock.Anything, "card-123").Return(nil)

	err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.NoError(t, err)
}

func TestDeleteCardUsecase_Execute_CardNotFound(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	logger := zap.NewNop()
	uc := New(cardRepo, logger)

	// Mock repository to return an error
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(nil, model.ErrCardNotFound)

	err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.Error(t, err)
	assert.Equal(t, model.ErrCardNotFound, err)
	// Repository delete should not be called
	cardRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteCardUsecase_Execute_UnauthorizedAccess(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	logger := zap.NewNop()
	uc := New(cardRepo, logger)

	existingCard := &model.Card{
		ID:     "card-123",
		UserID: "user-456", // Different user
		Title:  "My Card",
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.Error(t, err)
	assert.Equal(t, model.ErrForbidden, err)
	// Repository delete should not be called
	cardRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteCardUsecase_Execute_RepositoryError(t *testing.T) {
	cardRepo := mocks.NewCardRepository(t)
	logger := zap.NewNop()
	uc := New(cardRepo, logger)

	existingCard := &model.Card{
		ID:        "card-123",
		UserID:    "user-123",
		Title:     "My Card",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository to return the existing card
	cardRepo.EXPECT().GetByID(mock.Anything, "card-123").Return(existingCard, nil)

	// Mock repository to return an error on delete
	cardRepo.EXPECT().Delete(mock.Anything, "card-123").Return(assert.AnError)

	err := uc.Execute(context.Background(), "user-123", "card-123")

	assert.Error(t, err)
}
