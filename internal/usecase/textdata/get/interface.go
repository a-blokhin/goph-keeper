//go:generate mockery --name=GetTextDataUsecase --output=./mocks --outpkg=mocks --filename=get_text_data_usecase_mock.go --with-expecter

package textdataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetTextDataUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.TextData, error)
}
