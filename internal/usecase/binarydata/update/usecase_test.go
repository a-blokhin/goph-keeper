package update

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

func TestUpdateBinaryDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "user-123", Version: 1}
	updated := &model.BinaryData{ID: "bd-1", DataEncrypted: []byte("new data"), Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)
	encService.EXPECT().EncryptBinaryData(updated).Return(nil)
	repo.EXPECT().Update(mock.Anything, updated).Return(nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", updated.UserID)
	assert.Equal(t, int32(2), updated.Version)
}

func TestUpdateBinaryDataUsecase_Execute_TooLarge(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	updated := &model.BinaryData{
		ID:            "bd-1",
		DataEncrypted: make([]byte, 11*1024*1024),
		Version:       1,
	}

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrBinaryDataTooLarge)
}

func TestUpdateBinaryDataUsecase_Execute_Forbidden(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "other-user", Version: 1}
	updated := &model.BinaryData{ID: "bd-1", DataEncrypted: []byte("data"), Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrForbidden)
}

func TestUpdateBinaryDataUsecase_Execute_VersionConflict(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "user-123", Version: 2}
	updated := &model.BinaryData{ID: "bd-1", DataEncrypted: []byte("data"), Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrVersionConflict)
}

func TestUpdateBinaryDataUsecase_Execute_NotFound(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	updated := &model.BinaryData{ID: "bd-1", DataEncrypted: []byte("data"), Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(nil, errors.New("not found"))

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.Error(t, err)
}

func TestUpdateBinaryDataUsecase_Execute_EncryptionError(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "user-123", Version: 1}
	updated := &model.BinaryData{ID: "bd-1", DataEncrypted: []byte("data"), Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)
	encService.EXPECT().EncryptBinaryData(updated).Return(errors.New("encrypt failed"))

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.Error(t, err)
}
