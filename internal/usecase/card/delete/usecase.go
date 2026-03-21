package delete

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type deleteCardUsecase struct {
	cardRepo repository.CardRepository
	logger   *zap.Logger
}

func New(cardRepo repository.CardRepository, logger *zap.Logger) DeleteCardUsecase {
	return &deleteCardUsecase{
		cardRepo: cardRepo,
		logger:   logger,
	}
}

func (u *deleteCardUsecase) Execute(ctx context.Context, userID, cardID string) error {
	existing, err := u.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		u.logger.Error("failed to get card", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	err = u.cardRepo.Delete(ctx, cardID)
	if err != nil {
		u.logger.Error("failed to delete card", zap.Error(err))
		return err
	}

	return nil
}
