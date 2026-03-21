package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type listTextDataUsecase struct {
	textDataRepo repository.TextDataRepository
	logger       *zap.Logger
}

func New(textDataRepo repository.TextDataRepository, logger *zap.Logger) ListTextDataUsecase {
	return &listTextDataUsecase{
		textDataRepo: textDataRepo,
		logger:       logger,
	}
}

func (u *listTextDataUsecase) Execute(ctx context.Context, userID string) ([]*model.TextData, error) {
	textDataList, err := u.textDataRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get text data", zap.Error(err))
		return nil, err
	}

	return textDataList, nil
}
