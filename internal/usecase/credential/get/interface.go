package credentialget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type GetCredentialUsecase interface {
	Execute(ctx context.Context, userID, id string) (*model.Credential, error)
}
