package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type updateTextDataUsecase struct {
	textDataRepo repository.TextDataRepository
	logger       *zap.Logger
}

func New(textDataRepo repository.TextDataRepository, logger *zap.Logger) UpdateTextDataUsecase {
	return &updateTextDataUsecase{
		textDataRepo: textDataRepo,
		logger:       logger,
	}
}

func (u *updateTextDataUsecase) Execute(ctx context.Context, userID string, textData *model.TextData) error {
	existing, err := u.textDataRepo.GetByID(ctx, textData.ID)
	if err != nil {
		u.logger.Error("failed to get text data", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	if existing.Version != textData.Version {
		return model.ErrVersionConflict
	}

	textData.UserID = userID
	textData.Version = existing.Version + 1

	err = u.textDataRepo.Update(ctx, textData)
	if err != nil {
		u.logger.Error("failed to update text data", zap.Error(err))
		return err
	}

	return nil
}
