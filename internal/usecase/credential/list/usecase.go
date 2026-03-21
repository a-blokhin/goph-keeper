package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type listCredentialsUsecase struct {
	credentialRepo repository.CredentialRepository
	logger         *zap.Logger
}

func New(credentialRepo repository.CredentialRepository, logger *zap.Logger) ListCredentialsUsecase {
	return &listCredentialsUsecase{
		credentialRepo: credentialRepo,
		logger:         logger,
	}
}

func (u *listCredentialsUsecase) Execute(ctx context.Context, userID string) ([]*model.Credential, error) {
	credentials, err := u.credentialRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get credentials", zap.Error(err))
		return nil, err
	}

	return credentials, nil
}
