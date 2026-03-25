//go:generate mockery --name=GetCredentialUsecase --output=./mocks --outpkg=mocks --filename=get_credential_usecase_mock.go --with-expecter

package credentialget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetCredentialUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.Credential, error)
}
