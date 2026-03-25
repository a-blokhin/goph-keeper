//go:generate mockery --name=ListBinaryDataUsecase --output=./mocks --outpkg=mocks --filename=list_binary_data_usecase_mock.go --with-expecter

package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListBinaryDataUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.BinaryData, error)
}
