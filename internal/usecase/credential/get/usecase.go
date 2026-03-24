package credentialget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
)

type getCredentialUsecase struct {
	credentialRepo    repository.CredentialRepository
	encryptionService encryption.EncryptionService
}

func New(credentialRepo repository.CredentialRepository, encryptionService encryption.EncryptionService) GetCredentialUsecase {
	return &getCredentialUsecase{
		credentialRepo:    credentialRepo,
		encryptionService: encryptionService,
	}
}

func (u *getCredentialUsecase) Execute(ctx context.Context, userID, id string) (*model.Credential, error) {
	credential, err := u.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if credential.UserID != userID {
		return nil, model.ErrCredentialNotFound
	}

	if err := u.encryptionService.DecryptCredential(credential); err != nil {
		return nil, err
	}

	return credential, nil
}
