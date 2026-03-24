package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type updateCredentialUsecase struct {
	credentialRepo    repository.CredentialRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(credentialRepo repository.CredentialRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) UpdateCredentialUsecase {
	return &updateCredentialUsecase{
		credentialRepo:    credentialRepo,
		encryptionService: encryptionService,
		logger:            logger,
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

	if err := u.encryptionService.EncryptCredential(credential); err != nil {
		u.logger.Error("failed to encrypt credential", zap.Error(err))
		return err
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
