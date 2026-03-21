package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateBinaryDataUsecase interface {
	Execute(ctx context.Context, userID string, binaryData *model.BinaryData) error
}
