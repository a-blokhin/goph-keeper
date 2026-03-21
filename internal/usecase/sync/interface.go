package sync

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type SyncUsecase interface {
	Execute(ctx context.Context, userID string) (*SyncResponse, error)
}

type SyncResponse struct {
	Credentials []*model.Credential
	TextData    []*model.TextData
	BinaryData  []*model.BinaryData
	Cards       []*model.Card
	LastSync    time.Time
}
