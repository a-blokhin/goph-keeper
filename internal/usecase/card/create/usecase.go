package create

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type createCardUsecase struct {
	cardRepo          repository.CardRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(cardRepo repository.CardRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) CreateCardUsecase {
	return &createCardUsecase{
		cardRepo:          cardRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *createCardUsecase) Execute(ctx context.Context, userID string, card *model.Card) error {
	if err := u.encryptionService.EncryptCardData(card); err != nil {
		u.logger.Error("failed to encrypt card data", zap.Error(err))
		return err
	}

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
