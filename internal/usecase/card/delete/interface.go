package delete

import (
	"context"
)

type DeleteCardUsecase interface {
	Execute(ctx context.Context, userID, cardID string) error
}
