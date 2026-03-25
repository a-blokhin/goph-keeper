//go:generate mockery --name=DeleteCredentialUsecase --output=./mocks --outpkg=mocks --filename=delete_credential_usecase_mock.go --with-expecter

package delete

import (
	"context"
)

type DeleteCredentialUsecase interface {
	Execute(ctx context.Context, userID, credentialID string) error
}
