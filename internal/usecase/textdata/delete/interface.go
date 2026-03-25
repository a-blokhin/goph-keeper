//go:generate mockery --name=DeleteTextDataUsecase --output=./mocks --outpkg=mocks --filename=delete_text_data_usecase_mock.go --with-expecter

package delete

import (
	"context"
)

type DeleteTextDataUsecase interface {
	Execute(ctx context.Context, userID, textDataID string) error
}
