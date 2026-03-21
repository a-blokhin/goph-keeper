package cardget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
)

type getCardUsecase struct {
	cardRepo repository.CardRepository
}

func New(cardRepo repository.CardRepository) GetCardUsecase {
	return &getCardUsecase{
		cardRepo: cardRepo,
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

	return card, nil
}
