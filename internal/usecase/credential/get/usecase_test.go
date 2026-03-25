package credentialget

import (
	"context"
	"errors"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCredentialUsecase_Execute_Success(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(credRepo, encService)

	stored := &model.Credential{
		ID:                "cred-1",
		UserID:            "user-123",
		Title:             "My Login",
		Login:             "user",
		PasswordEncrypted: "encrypted-pw",
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(stored, nil)
	encService.EXPECT().DecryptCredential(stored).Return(nil)

	result, err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.NoError(t, err)
	assert.Equal(t, stored, result)
}

func TestGetCredentialUsecase_Execute_WrongUser(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(credRepo, encService)

	stored := &model.Credential{
		ID:     "cred-1",
		UserID: "other-user",
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(stored, nil)

	result, err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.ErrorIs(t, err, model.ErrCredentialNotFound)
	assert.Nil(t, result)
}

func TestGetCredentialUsecase_Execute_NotFound(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(credRepo, encService)

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(nil, errors.New("not found"))

	result, err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetCredentialUsecase_Execute_DecryptError(t *testing.T) {
	credRepo := mocks.NewCredentialRepository(t)
	encService := encmocks.NewEncryptionService(t)

	uc := New(credRepo, encService)

	stored := &model.Credential{
		ID:                "cred-1",
		UserID:            "user-123",
		PasswordEncrypted: "corrupted",
	}

	credRepo.EXPECT().GetByID(mock.Anything, "cred-1").Return(stored, nil)
	encService.EXPECT().DecryptCredential(stored).Return(errors.New("decrypt failed"))

	result, err := uc.Execute(context.Background(), "user-123", "cred-1")

	assert.Error(t, err)
	assert.Nil(t, result)
}
