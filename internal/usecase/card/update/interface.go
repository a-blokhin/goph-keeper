//go:generate mockery --name=UpdateCardUsecase --output=./mocks --outpkg=mocks --filename=update_card_usecase_mock.go --with-expecter

package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateCardUsecase interface {
	Execute(ctx context.Context, userID string, card *model.Card) error
}
