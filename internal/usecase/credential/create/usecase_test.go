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

func TestCreateCredentialUsecase_Execute_Success(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	cred := &model.Credential{
		Title:             "My Login",
		Login:             "user",
		PasswordEncrypted: "secret",
	}

	encService.EXPECT().EncryptCredential(cred).Return(nil)
	credRepo.EXPECT().Create(mock.Anything, cred).Return(nil)

	err := uc.Execute(context.Background(), "user-123", cred)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", cred.UserID)
	assert.Equal(t, int32(1), cred.Version)
	assert.False(t, cred.CreatedAt.IsZero())
	assert.False(t, cred.UpdatedAt.IsZero())
}

func TestCreateCredentialUsecase_Execute_EncryptionError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	cred := &model.Credential{
		Title:             "My Login",
		Login:             "user",
		PasswordEncrypted: "secret",
	}

	encService.EXPECT().EncryptCredential(cred).Return(errors.New("encryption failed"))

	err := uc.Execute(context.Background(), "user-123", cred)

	assert.Error(t, err)
	assert.Equal(t, "encryption failed", err.Error())
}

func TestCreateCredentialUsecase_Execute_RepoError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	cred := &model.Credential{
		Title:             "My Login",
		Login:             "user",
		PasswordEncrypted: "secret",
	}

	encService.EXPECT().EncryptCredential(cred).Return(nil)
	credRepo.EXPECT().Create(mock.Anything, cred).Return(errors.New("db error"))

	err := uc.Execute(context.Background(), "user-123", cred)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
