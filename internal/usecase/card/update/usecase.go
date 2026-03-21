package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type updateCardUsecase struct {
	cardRepo repository.CardRepository
	logger   *zap.Logger
}

func New(cardRepo repository.CardRepository, logger *zap.Logger) UpdateCardUsecase {
	return &updateCardUsecase{
		cardRepo: cardRepo,
		logger:   logger,
	}
}

func (u *updateCardUsecase) Execute(ctx context.Context, userID string, card *model.Card) error {
	existing, err := u.cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		u.logger.Error("failed to get card", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	if existing.Version != card.Version {
		return model.ErrVersionConflict
	}

	card.UserID = userID
	card.Version = existing.Version + 1

	err = u.cardRepo.Update(ctx, card)
	if err != nil {
		u.logger.Error("failed to update card", zap.Error(err))
		return err
	}

	return nil
}
