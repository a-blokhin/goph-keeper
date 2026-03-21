package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateBinaryDataUsecase interface {
	Execute(ctx context.Context, userID string, binaryData *model.BinaryData) error
}
