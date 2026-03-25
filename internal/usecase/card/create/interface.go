//go:generate mockery --name=CreateCardUsecase --output=./mocks --outpkg=mocks --filename=create_card_usecase_mock.go --with-expecter

package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateCardUsecase interface {
	Execute(ctx context.Context, userID string, card *model.Card) error
}
