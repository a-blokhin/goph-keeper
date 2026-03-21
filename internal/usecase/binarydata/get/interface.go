package binarydataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetBinaryDataUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.BinaryData, error)
}
