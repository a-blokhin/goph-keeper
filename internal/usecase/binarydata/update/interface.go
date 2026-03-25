//go:generate mockery --name=UpdateBinaryDataUsecase --output=./mocks --outpkg=mocks --filename=update_binary_data_usecase_mock.go --with-expecter

package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateBinaryDataUsecase interface {
	Execute(ctx context.Context, userID string, binaryData *model.BinaryData) error
}
