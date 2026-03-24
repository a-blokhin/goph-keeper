package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type listCredentialsUsecase struct {
	credentialRepo    repository.CredentialRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(credentialRepo repository.CredentialRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) ListCredentialsUsecase {
	return &listCredentialsUsecase{
		credentialRepo:    credentialRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *listCredentialsUsecase) Execute(ctx context.Context, userID string) ([]*model.Credential, error) {
	credentials, err := u.credentialRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get credentials", zap.Error(err))
		return nil, err
	}

	for _, credential := range credentials {
		if err := u.encryptionService.DecryptCredential(credential); err != nil {
			u.logger.Error("failed to decrypt credential", zap.Error(err))
			return nil, err
		}
	}

	return credentials, nil
}
