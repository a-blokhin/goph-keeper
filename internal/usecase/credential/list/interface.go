//go:generate mockery --name=ListCredentialsUsecase --output=./mocks --outpkg=mocks --filename=list_credentials_usecase_mock.go --with-expecter

package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListCredentialsUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.Credential, error)
}
