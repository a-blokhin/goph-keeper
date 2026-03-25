package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type listCardsUsecase struct {
	cardRepo          repository.CardRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(cardRepo repository.CardRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) ListCardsUsecase {
	return &listCardsUsecase{
		cardRepo:          cardRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *listCardsUsecase) Execute(ctx context.Context, userID string) ([]*model.Card, error) {
	cards, err := u.cardRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get cards", zap.Error(err))
		return nil, err
	}

	for _, card := range cards {
		if err := u.encryptionService.DecryptCardData(card); err != nil {
			u.logger.Error("failed to decrypt card data", zap.Error(err))
			return nil, err
		}
	}

	return cards, nil
}
