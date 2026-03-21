package delete

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type deleteTextDataUsecase struct {
	textDataRepo repository.TextDataRepository
	logger       *zap.Logger
}

func New(textDataRepo repository.TextDataRepository, logger *zap.Logger) DeleteTextDataUsecase {
	return &deleteTextDataUsecase{
		textDataRepo: textDataRepo,
		logger:       logger,
	}
}

func (u *deleteTextDataUsecase) Execute(ctx context.Context, userID, textDataID string) error {
	existing, err := u.textDataRepo.GetByID(ctx, textDataID)
	if err != nil {
		u.logger.Error("failed to get text data", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	err = u.textDataRepo.Delete(ctx, textDataID)
	if err != nil {
		u.logger.Error("failed to delete text data", zap.Error(err))
		return err
	}

	return nil
}
