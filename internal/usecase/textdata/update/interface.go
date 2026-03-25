//go:generate mockery --name=UpdateTextDataUsecase --output=./mocks --outpkg=mocks --filename=update_text_data_usecase_mock.go --with-expecter

package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateTextDataUsecase interface {
	Execute(ctx context.Context, userID string, textData *model.TextData) error
}
