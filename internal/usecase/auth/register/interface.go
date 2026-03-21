package register

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type RegisterUsecase interface {
	Execute(ctx context.Context, email, password string) (*model.User, string, error)
}
