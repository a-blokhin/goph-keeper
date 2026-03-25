package textdataget

import (
	"context"
	"errors"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTextDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(repo, encService)

	stored := &model.TextData{
		ID:            "td-1",
		UserID:        "user-123",
		DataEncrypted: "encrypted",
	}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(stored, nil)
	encService.EXPECT().DecryptTextData(stored).Return(nil)

	result, err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.NoError(t, err)
	assert.Equal(t, stored, result)
}

func TestGetTextDataUsecase_Execute_WrongUser(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(repo, encService)

	stored := &model.TextData{
		ID:     "td-1",
		UserID: "other-user",
	}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(stored, nil)

	result, err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.ErrorIs(t, err, model.ErrTextDataNotFound)
	assert.Nil(t, result)
}

func TestGetTextDataUsecase_Execute_NotFound(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(repo, encService)

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(nil, errors.New("not found"))

	result, err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTextDataUsecase_Execute_DecryptError(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(repo, encService)

	stored := &model.TextData{
		ID:            "td-1",
		UserID:        "user-123",
		DataEncrypted: "corrupted",
	}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(stored, nil)
	encService.EXPECT().DecryptTextData(stored).Return(errors.New("decrypt failed"))

	result, err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}
