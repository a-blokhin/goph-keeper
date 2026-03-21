package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListCardsUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.Card, error)
}
