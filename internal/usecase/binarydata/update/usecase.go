package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

const maxBinaryDataSize = 10 * 1024 * 1024

type updateBinaryDataUsecase struct {
	binaryDataRepo    repository.BinaryDataRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(binaryDataRepo repository.BinaryDataRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) UpdateBinaryDataUsecase {
	return &updateBinaryDataUsecase{
		binaryDataRepo:    binaryDataRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *updateBinaryDataUsecase) Execute(ctx context.Context, userID string, binaryData *model.BinaryData) error {
	if len(binaryData.DataEncrypted) > maxBinaryDataSize {
		return model.ErrBinaryDataTooLarge
	}

	existing, err := u.binaryDataRepo.GetByID(ctx, binaryData.ID)
	if err != nil {
		u.logger.Error("failed to get binary data", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	if existing.Version != binaryData.Version {
		return model.ErrVersionConflict
	}

	if err := u.encryptionService.EncryptBinaryData(binaryData); err != nil {
		u.logger.Error("failed to encrypt binary data", zap.Error(err))
		return err
	}

	binaryData.UserID = userID
	binaryData.Version = existing.Version + 1

	err = u.binaryDataRepo.Update(ctx, binaryData)
	if err != nil {
		u.logger.Error("failed to update binary data", zap.Error(err))
		return err
	}

	return nil
}
