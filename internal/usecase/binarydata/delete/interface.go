//go:generate mockery --name=DeleteBinaryDataUsecase --output=./mocks --outpkg=mocks --filename=delete_binary_data_usecase_mock.go --with-expecter

package delete

import (
	"context"
)

type DeleteBinaryDataUsecase interface {
	Execute(ctx context.Context, userID, binaryDataID string) error
}
