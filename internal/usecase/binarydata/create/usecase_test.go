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

func TestCreateBinaryDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	bd := &model.BinaryData{
		Title:         "My File",
		DataEncrypted: []byte("file content"),
	}

	encService.EXPECT().EncryptBinaryData(bd).Return(nil)
	repo.EXPECT().Create(mock.Anything, bd).Return(nil)

	err := uc.Execute(context.Background(), "user-123", bd)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", bd.UserID)
	assert.Equal(t, int32(1), bd.Version)
	assert.False(t, bd.CreatedAt.IsZero())
	assert.False(t, bd.UpdatedAt.IsZero())
}

func TestCreateBinaryDataUsecase_Execute_TooLarge(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	bd := &model.BinaryData{
		Title:         "Huge File",
		DataEncrypted: make([]byte, 11*1024*1024), // 11MB > 10MB limit
	}

	err := uc.Execute(context.Background(), "user-123", bd)

	assert.ErrorIs(t, err, model.ErrBinaryDataTooLarge)
}

func TestCreateBinaryDataUsecase_Execute_ExactlyAtLimit(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	bd := &model.BinaryData{
		Title:         "Exact Limit File",
		DataEncrypted: make([]byte, 10*1024*1024), // exactly 10MB
	}

	encService.EXPECT().EncryptBinaryData(bd).Return(nil)
	repo.EXPECT().Create(mock.Anything, bd).Return(nil)

	err := uc.Execute(context.Background(), "user-123", bd)

	assert.NoError(t, err)
}

func TestCreateBinaryDataUsecase_Execute_EncryptionError(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	bd := &model.BinaryData{DataEncrypted: []byte("data")}

	encService.EXPECT().EncryptBinaryData(bd).Return(errors.New("encryption failed"))

	err := uc.Execute(context.Background(), "user-123", bd)

	assert.Error(t, err)
}

func TestCreateBinaryDataUsecase_Execute_RepoError(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	bd := &model.BinaryData{DataEncrypted: []byte("data")}

	encService.EXPECT().EncryptBinaryData(bd).Return(nil)
	repo.EXPECT().Create(mock.Anything, bd).Return(errors.New("db error"))

	err := uc.Execute(context.Background(), "user-123", bd)

	assert.Error(t, err)
}
