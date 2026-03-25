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

func TestUpdateCredentialUsecase_Execute_Success(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	existing := &model.Credential{
		ID:      "cred-1",
		UserID:  "user-123",
		Version: 1,
	}

	updated := &model.Credential{
		ID:                "cred-1",
		Title:             "Updated Title",
		Login:             "newlogin",
		PasswordEncrypted: "newsecret",
		Version:           1,
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)
	encService.EXPECT().EncryptCredential(updated).Return(nil)
	credRepo.EXPECT().Update(mock.Anything, updated).Return(nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", updated.UserID)
	assert.Equal(t, int32(2), updated.Version)
}

func TestUpdateCredentialUsecase_Execute_Forbidden(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	existing := &model.Credential{
		ID:      "cred-1",
		UserID:  "other-user",
		Version: 1,
	}

	updated := &model.Credential{
		ID:      "cred-1",
		Version: 1,
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrForbidden)
}

func TestUpdateCredentialUsecase_Execute_VersionConflict(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	existing := &model.Credential{
		ID:      "cred-1",
		UserID:  "user-123",
		Version: 2,
	}

	updated := &model.Credential{
		ID:      "cred-1",
		Version: 1,
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.ErrorIs(t, err, model.ErrVersionConflict)
}

func TestUpdateCredentialUsecase_Execute_NotFound(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	updated := &model.Credential{
		ID:      "cred-1",
		Version: 1,
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(nil, errors.New("not found"))

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.Error(t, err)
}

func TestUpdateCredentialUsecase_Execute_EncryptionError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	existing := &model.Credential{
		ID:      "cred-1",
		UserID:  "user-123",
		Version: 1,
	}

	updated := &model.Credential{
		ID:      "cred-1",
		Version: 1,
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(existing, nil)
	encService.EXPECT().EncryptCredential(updated).Return(errors.New("encrypt failed"))

	err := uc.Execute(context.Background(), "user-123", updated)

	assert.Error(t, err)
	assert.Equal(t, "encrypt failed", err.Error())
}
