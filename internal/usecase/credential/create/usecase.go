package create

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type createCredentialUsecase struct {
	credentialRepo repository.CredentialRepository
	logger         *zap.Logger
}

func New(credentialRepo repository.CredentialRepository, logger *zap.Logger) CreateCredentialUsecase {
	return &createCredentialUsecase{
		credentialRepo: credentialRepo,
		logger:         logger,
	}
}

func (u *createCredentialUsecase) Execute(ctx context.Context, userID string, credential *model.Credential) error {
	credential.UserID = userID
	credential.Version = 1
	now := time.Now()
	credential.CreatedAt = now
	credential.UpdatedAt = now

	err := u.credentialRepo.Create(ctx, credential)
	if err != nil {
		u.logger.Error("failed to create credential", zap.Error(err))
		return err
	}

	return nil
}
