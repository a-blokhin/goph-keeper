//go:generate mockery --name=GetBinaryDataUsecase --output=./mocks --outpkg=mocks --filename=get_binary_data_usecase_mock.go --with-expecter

package binarydataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetBinaryDataUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.BinaryData, error)
}
