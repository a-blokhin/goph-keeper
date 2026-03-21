package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateCardUsecase interface {
	Execute(ctx context.Context, userID string, card *model.Card) error
}
