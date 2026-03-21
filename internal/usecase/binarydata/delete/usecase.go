package delete

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type deleteBinaryDataUsecase struct {
	binaryDataRepo repository.BinaryDataRepository
	logger         *zap.Logger
}

func New(binaryDataRepo repository.BinaryDataRepository, logger *zap.Logger) DeleteBinaryDataUsecase {
	return &deleteBinaryDataUsecase{
		binaryDataRepo: binaryDataRepo,
		logger:         logger,
	}
}

func (u *deleteBinaryDataUsecase) Execute(ctx context.Context, userID, binaryDataID string) error {
	existing, err := u.binaryDataRepo.GetByID(ctx, binaryDataID)
	if err != nil {
		u.logger.Error("failed to get binary data", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	err = u.binaryDataRepo.Delete(ctx, binaryDataID)
	if err != nil {
		u.logger.Error("failed to delete binary data", zap.Error(err))
		return err
	}

	return nil
}
