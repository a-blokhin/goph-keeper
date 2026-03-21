package cardget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetCardUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.Card, error)
}
