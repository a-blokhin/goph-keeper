package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestDeleteBinaryDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "user-123"}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)
	repo.EXPECT().Delete(mock.Anything, "bd-1").Return(nil)

	err := uc.Execute(context.Background(), "user-123", "bd-1")

	assert.NoError(t, err)
}

func TestDeleteBinaryDataUsecase_Execute_Forbidden(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "other-user"}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", "bd-1")

	assert.ErrorIs(t, err, model.ErrForbidden)
}

func TestDeleteBinaryDataUsecase_Execute_NotFound(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(nil, errors.New("not found"))

	err := uc.Execute(context.Background(), "user-123", "bd-1")

	assert.Error(t, err)
}

func TestDeleteBinaryDataUsecase_Execute_DeleteError(t *testing.T) {
	repo := mocks.NewBinaryDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	existing := &model.BinaryData{ID: "bd-1", UserID: "user-123"}

	repo.EXPECT().GetByID(mock.Anything, "bd-1").Return(existing, nil)
	repo.EXPECT().Delete(mock.Anything, "bd-1").Return(errors.New("db error"))

	err := uc.Execute(context.Background(), "user-123", "bd-1")

	assert.Error(t, err)
}
