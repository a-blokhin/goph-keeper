//go:generate mockery --name=CreateCredentialUsecase --output=./mocks --outpkg=mocks --filename=create_credential_usecase_mock.go --with-expecter

package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateCredentialUsecase interface {
	Execute(ctx context.Context, userID string, credential *model.Credential) error
}
