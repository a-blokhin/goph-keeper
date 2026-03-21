package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateTextDataUsecase interface {
	Execute(ctx context.Context, userID string, textData *model.TextData) error
}
