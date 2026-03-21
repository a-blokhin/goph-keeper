package create

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type createCardUsecase struct {
	cardRepo repository.CardRepository
	logger   *zap.Logger
}

func New(cardRepo repository.CardRepository, logger *zap.Logger) CreateCardUsecase {
	return &createCardUsecase{
		cardRepo: cardRepo,
		logger:   logger,
	}
}

func (u *createCardUsecase) Execute(ctx context.Context, userID string, card *model.Card) error {
	card.UserID = userID
	card.Version = 1
	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	err := u.cardRepo.Create(ctx, card)
	if err != nil {
		u.logger.Error("failed to create card", zap.Error(err))
		return err
	}

	return nil
}
