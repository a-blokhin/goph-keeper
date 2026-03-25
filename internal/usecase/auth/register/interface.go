//go:generate mockery --name=RegisterUsecase --output=./mocks --outpkg=mocks --filename=register_usecase_mock.go --with-expecter

package register

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type RegisterUsecase interface {
	Execute(ctx context.Context, email, password string) (*model.User, string, error)
}
