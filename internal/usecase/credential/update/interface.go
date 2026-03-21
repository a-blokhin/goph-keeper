package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UpdateCredentialUsecase interface {
	Execute(ctx context.Context, userID string, credential *model.Credential) error
}
