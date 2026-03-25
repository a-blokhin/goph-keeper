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

func TestListCredentialsUsecase_Execute_Success(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	creds := []*model.Credential{
		{ID: "cred-1", UserID: "user-123", PasswordEncrypted: "enc1"},
		{ID: "cred-2", UserID: "user-123", PasswordEncrypted: "enc2"},
	}

	credRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(creds, nil)
	encService.EXPECT().DecryptCredential(creds[0]).Return(nil)
	encService.EXPECT().DecryptCredential(creds[1]).Return(nil)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListCredentialsUsecase_Execute_Empty(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	credRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return([]*model.Credential{}, nil)

	result, err := uc.Execute(context.Background(), "user-123")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestListCredentialsUsecase_Execute_RepoError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	credRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(nil, errors.New("db error"))

	result, err := uc.Execute(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListCredentialsUsecase_Execute_DecryptError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)
	logger := zap.NewNop()

	uc := New(credRepo, encService, logger)

	creds := []*model.Credential{
		{ID: "cred-1", UserID: "user-123", PasswordEncrypted: "enc1"},
	}

	credRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(creds, nil)
	encService.EXPECT().DecryptCredential(creds[0]).Return(errors.New("decrypt failed"))

	result, err := uc.Execute(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}
