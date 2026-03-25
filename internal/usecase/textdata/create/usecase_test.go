package create

import (
	"context"
	"errors"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateTextDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	td := &model.TextData{
		Title:         "My Note",
		DataEncrypted: "secret text",
	}

	encService.EXPECT().EncryptTextData(td).Return(nil)
	repo.EXPECT().Create(mock.Anything, td).Return(nil)

	err := uc.Execute(context.Background(), "user-123", td)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", td.UserID)
	assert.Equal(t, int32(1), td.Version)
	assert.False(t, td.CreatedAt.IsZero())
	assert.False(t, td.UpdatedAt.IsZero())
}

func TestCreateTextDataUsecase_Execute_EncryptionError(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	td := &model.TextData{DataEncrypted: "text"}

	encService.EXPECT().EncryptTextData(td).Return(errors.New("encryption failed"))

	err := uc.Execute(context.Background(), "user-123", td)

	assert.Error(t, err)
}

func TestCreateTextDataUsecase_Execute_RepoError(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	td := &model.TextData{DataEncrypted: "text"}

	encService.EXPECT().EncryptTextData(td).Return(nil)
	repo.EXPECT().Create(mock.Anything, td).Return(errors.New("db error"))

	err := uc.Execute(context.Background(), "user-123", td)

	assert.Error(t, err)
}
