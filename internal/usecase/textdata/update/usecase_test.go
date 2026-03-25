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

func TestUpdateTextDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.TextData{ID: "td-1", UserID: "user-123", Version: 1}
	updated := &model.TextData{ID: "td-1", DataEncrypted: "new text", Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)
	encService.EXPECT().EncryptTextData(updated).Return(nil)
	repo.EXPECT().Update(mock.Anything, updated).Return(nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", updated.UserID)
	assert.Equal(t, int32(2), updated.Version)
}

func TestUpdateTextDataUsecase_Execute_Forbidden(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.TextData{ID: "td-1", UserID: "other-user", Version: 1}
	updated := &model.TextData{ID: "td-1", Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrForbidden)
}

func TestUpdateTextDataUsecase_Execute_VersionConflict(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.TextData{ID: "td-1", UserID: "user-123", Version: 2}
	updated := &model.TextData{ID: "td-1", Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrVersionConflict)
}

func TestUpdateTextDataUsecase_Execute_NotFound(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	updated := &model.TextData{ID: "td-1", Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(nil, errors.New("not found"))

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.Error(t, err)
}

func TestUpdateTextDataUsecase_Execute_EncryptionError(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(repo, encService, logger)

	existing := &model.TextData{ID: "td-1", UserID: "user-123", Version: 1}
	updated := &model.TextData{ID: "td-1", Version: 1}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)
	encService.EXPECT().EncryptTextData(updated).Return(errors.New("encrypt failed"))

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.Error(t, err)
}
