//go:generate mockery --name=DeleteCardUsecase --output=./mocks --outpkg=mocks --filename=delete_card_usecase_mock.go --with-expecter

package delete

import (
	"context"
)

type DeleteCardUsecase interface {
	Execute(ctx context.Context, userID, cardID string) error
}
