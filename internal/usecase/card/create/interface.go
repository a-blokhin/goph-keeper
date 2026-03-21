package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateCardUsecase interface {
	Execute(ctx context.Context, userID string, card *model.Card) error
}
