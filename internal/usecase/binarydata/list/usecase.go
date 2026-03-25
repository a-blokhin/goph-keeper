package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type listBinaryDataUsecase struct {
	binaryDataRepo    repository.BinaryDataRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(binaryDataRepo repository.BinaryDataRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) ListBinaryDataUsecase {
	return &listBinaryDataUsecase{
		binaryDataRepo:    binaryDataRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *listBinaryDataUsecase) Execute(ctx context.Context, userID string) ([]*model.BinaryData, error) {
	binaryDataList, err := u.binaryDataRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get binary data", zap.Error(err))
		return nil, err
	}

	for _, binaryData := range binaryDataList {
		if err := u.encryptionService.DecryptBinaryData(binaryData); err != nil {
			u.logger.Error("failed to decrypt binary data", zap.Error(err))
			return nil, err
		}
	}

	return binaryDataList, nil
}
