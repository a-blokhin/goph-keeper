package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type updateCredentialUsecase struct {
	credentialRepo repository.CredentialRepository
	logger         *zap.Logger
}

func New(credentialRepo repository.CredentialRepository, logger *zap.Logger) UpdateCredentialUsecase {
	return &updateCredentialUsecase{
		credentialRepo: credentialRepo,
		logger:         logger,
	}
}

func (u *updateCredentialUsecase) Execute(ctx context.Context, userID string, credential *model.Credential) error {
	existing, err := u.credentialRepo.GetByID(ctx, credential.ID)
	if err != nil {
		u.logger.Error("failed to get credential", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	if existing.Version != credential.Version {
		return model.ErrVersionConflict
	}

	credential.UserID = userID
	credential.Version = existing.Version + 1

	err = u.credentialRepo.Update(ctx, credential)
	if err != nil {
		u.logger.Error("failed to update credential", zap.Error(err))
		return err
	}

	return nil
}
