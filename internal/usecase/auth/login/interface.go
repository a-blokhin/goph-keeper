//go:generate mockery --name=LoginUsecase --output=./mocks --outpkg=mocks --filename=login_usecase_mock.go --with-expecter

package login

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type LoginUsecase interface {
	Execute(ctx context.Context, email, password string) (*model.User, string, error)
}
