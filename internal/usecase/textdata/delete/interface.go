package delete

import (
	"context"
)

type DeleteTextDataUsecase interface {
	Execute(ctx context.Context, userID, textDataID string) error
}
