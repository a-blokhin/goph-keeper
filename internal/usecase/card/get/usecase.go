package cardget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
)

type getCardUsecase struct {
	cardRepo          repository.CardRepository
	encryptionService encryption.EncryptionService
}

func New(cardRepo repository.CardRepository, encryptionService encryption.EncryptionService) GetCardUsecase {
	return &getCardUsecase{
		cardRepo:          cardRepo,
		encryptionService: encryptionService,
	}
}

func (u *getCardUsecase) Execute(ctx context.Context, userID, id string) (*model.Card, error) {
	card, err := u.cardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if card.UserID != userID {
		return nil, model.ErrCardNotFound
	}

	if err := u.encryptionService.DecryptCardData(card); err != nil {
		return nil, err
	}

	return card, nil
}
