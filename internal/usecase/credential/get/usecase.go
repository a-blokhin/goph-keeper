package credentialget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
)

type getCredentialUsecase struct {
	credentialRepo repository.CredentialRepository
}

func New(credentialRepo repository.CredentialRepository) GetCredentialUsecase {
	return &getCredentialUsecase{
		credentialRepo: credentialRepo,
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

	return credential, nil
}
