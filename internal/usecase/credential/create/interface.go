package create

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type CreateCredentialUsecase interface {
	Execute(ctx context.Context, userID string, credential *model.Credential) error
}
