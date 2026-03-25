//go:generate mockery --name=CreateBinaryDataUsecase --output=./mocks --outpkg=mocks --filename=create_binary_data_usecase_mock.go --with-expecter

package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateBinaryDataUsecase interface {
	Execute(ctx context.Context, userID string, binaryData *model.BinaryData) error
}
