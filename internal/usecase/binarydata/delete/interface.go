package delete

import (
	"context"
)

type DeleteBinaryDataUsecase interface {
	Execute(ctx context.Context, userID, binaryDataID string) error
}
