//go:generate mockery --name=ListCardsUsecase --output=./mocks --outpkg=mocks --filename=list_cards_usecase_mock.go --with-expecter

package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListCardsUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.Card, error)
}
