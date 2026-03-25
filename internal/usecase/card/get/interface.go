//go:generate mockery --name=GetCardUsecase --output=./mocks --outpkg=mocks --filename=get_card_usecase_mock.go --with-expecter

package cardget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetCardUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.Card, error)
}
