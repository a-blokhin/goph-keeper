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

func TestDeleteCredentialUsecase_Execute_Success(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	logger := zap.NewNop()

	uc := New(credRepo, logger)

	existing := &model.Credential{
		ID:     "cred-1",
		UserID: "user-123",
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)
	credRepo.EXPECT().Delete(mock.Anything, "cred-1").Return(nil)

	err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.NoError(t, err)
}

func TestDeleteCredentialUsecase_Execute_Forbidden(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	logger := zap.NewNop()

	uc := New(credRepo, logger)

	existing := &model.Credential{
		ID:     "cred-1",
		UserID: "other-user",
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.ErrorIs(t, err, model.ErrForbidden)
}

func TestDeleteCredentialUsecase_Execute_NotFound(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	logger := zap.NewNop()

	uc := New(credRepo, logger)

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(nil, errors.New("not found"))

	err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.Error(t, err)
}

func TestDeleteCredentialUsecase_Execute_DeleteError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	logger := zap.NewNop()

	uc := New(credRepo, logger)

	existing := &model.Credential{
		ID:     "cred-1",
		UserID: "user-123",
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)
	credRepo.EXPECT().Delete(mock.Anything, "cred-1").Return(errors.New("db error"))

	err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
