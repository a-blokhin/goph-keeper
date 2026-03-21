package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListTextDataUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.TextData, error)
}
