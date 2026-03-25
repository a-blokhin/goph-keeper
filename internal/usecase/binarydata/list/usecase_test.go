package list

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

func TestListBinaryDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	items := []*model.BinaryData{
		{ID: "bd-1", UserID: "user-123", DataEncrypted: []byte("enc1")},
		{ID: "bd-2", UserID: "user-123", DataEncrypted: []byte("enc2")},
	}

	repo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(items, nil)
	encService.EXPECT().DecryptBinaryData(items[0]).Return(nil)
	encService.EXPECT().DecryptBinaryData(items[1]).Return(nil)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListBinaryDataUsecase_Execute_Empty(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	repo.EXPECT().GetByUserID(mock.Anything, "user-123").Return([]*model.BinaryData{}, nil)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestListBinaryDataUsecase_Execute_RepoError(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	repo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(nil, errors.New("db error"))

	result, err := uc.Execute(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListBinaryDataUsecase_Execute_DecryptError(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	items := []*model.BinaryData{
		{ID: "bd-1", UserID: "user-123", DataEncrypted: []byte("enc1")},
	}

	repo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(items, nil)
	encService.EXPECT().DecryptBinaryData(items[0]).Return(errors.New("decrypt failed"))

	result, err := uc.Execute(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}
