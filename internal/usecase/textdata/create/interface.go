//go:generate mockery --name=CreateTextDataUsecase --output=./mocks --outpkg=mocks --filename=create_text_data_usecase_mock.go --with-expecter

package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateTextDataUsecase interface {
	Execute(ctx context.Context, userID string, textData *model.TextData) error
}
