package create

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

const maxBinaryDataSize = 10 * 1024 * 1024

type createBinaryDataUsecase struct {
	binaryDataRepo repository.BinaryDataRepository
	logger         *zap.Logger
}

func New(binaryDataRepo repository.BinaryDataRepository, logger *zap.Logger) CreateBinaryDataUsecase {
	return &createBinaryDataUsecase{
		binaryDataRepo: binaryDataRepo,
		logger:         logger,
	}
}

func (u *createBinaryDataUsecase) Execute(ctx context.Context, userID string, binaryData *model.BinaryData) error {
	if len(binaryData.DataEncrypted) > maxBinaryDataSize {
		return model.ErrBinaryDataTooLarge
	}

	binaryData.UserID = userID
	binaryData.Version = 1
	now := time.Now()
	binaryData.CreatedAt = now
	binaryData.UpdatedAt = now

	err := u.binaryDataRepo.Create(ctx, binaryData)
	if err != nil {
		u.logger.Error("failed to create binary data", zap.Error(err))
		return err
	}

	return nil
}
