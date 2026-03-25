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

func TestDeleteTextDataUsecase_Execute_Success(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	existing := &model.TextData{ID: "td-1", UserID: "user-123"}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)
	repo.EXPECT().Delete(mock.Anything, "td-1").Return(nil)

	err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.NoError(t, err)
}

func TestDeleteTextDataUsecase_Execute_Forbidden(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	existing := &model.TextData{ID: "td-1", UserID: "other-user"}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.ErrorIs(t, err, model.ErrForbidden)
}

func TestDeleteTextDataUsecase_Execute_NotFound(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(nil, errors.New("not found"))

	err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.Error(t, err)
}

func TestDeleteTextDataUsecase_Execute_DeleteError(t *testing.T) {
	repo := mocks.NewTextDataRepository(t)
	logger := zap.NewNop()

	uc := New(repo, logger)

	existing := &model.TextData{ID: "td-1", UserID: "user-123"}

	repo.EXPECT().GetByID(mock.Anything, "td-1").Return(existing, nil)
	repo.EXPECT().Delete(mock.Anything, "td-1").Return(errors.New("db error"))

	err := uc.Execute(context.Background(), "user-123", "td-1")

	assert.Error(t, err)
}
