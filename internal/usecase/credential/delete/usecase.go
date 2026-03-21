package delete

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type deleteCredentialUsecase struct {
	credentialRepo repository.CredentialRepository
	logger         *zap.Logger
}

func New(credentialRepo repository.CredentialRepository, logger *zap.Logger) DeleteCredentialUsecase {
	return &deleteCredentialUsecase{
		credentialRepo: credentialRepo,
		logger:         logger,
	}
}

func (u *deleteCredentialUsecase) Execute(ctx context.Context, userID, credentialID string) error {
	existing, err := u.credentialRepo.GetByID(ctx, credentialID)
	if err != nil {
		u.logger.Error("failed to get credential", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	err = u.credentialRepo.Delete(ctx, credentialID)
	if err != nil {
		u.logger.Error("failed to delete credential", zap.Error(err))
		return err
	}

	return nil
}
