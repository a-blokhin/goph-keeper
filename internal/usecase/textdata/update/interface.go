package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateTextDataUsecase interface {
	Execute(ctx context.Context, userID string, textData *model.TextData) error
}
