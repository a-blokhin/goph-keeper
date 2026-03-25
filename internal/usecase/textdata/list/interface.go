//go:generate mockery --name=ListTextDataUsecase --output=./mocks --outpkg=mocks --filename=list_text_data_usecase_mock.go --with-expecter

package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListTextDataUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.TextData, error)
}
