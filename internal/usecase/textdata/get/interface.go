package textdataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetTextDataUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.TextData, error)
}
