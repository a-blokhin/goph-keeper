package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type listCardsUsecase struct {
	cardRepo repository.CardRepository
	logger   *zap.Logger
}

func New(cardRepo repository.CardRepository, logger *zap.Logger) ListCardsUsecase {
	return &listCardsUsecase{
		cardRepo: cardRepo,
		logger:   logger,
	}
}

func (u *listCardsUsecase) Execute(ctx context.Context, userID string) ([]*model.Card, error) {
	cards, err := u.cardRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get cards", zap.Error(err))
		return nil, err
	}

	return cards, nil
}
