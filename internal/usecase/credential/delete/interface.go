package delete

import (
	"context"
)

type DeleteCredentialUsecase interface {
	Execute(ctx context.Context, userID, credentialID string) error
}
