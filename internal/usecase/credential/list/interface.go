package list

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type ListCredentialsUsecase interface {
	Execute(ctx context.Context, userID string) ([]*model.Credential, error)
}
