package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListBinaryDataUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.BinaryData, error)
}
