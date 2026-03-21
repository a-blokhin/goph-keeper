package sync

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type syncUsecase struct {
	credentialRepo repository.CredentialRepository
	textDataRepo   repository.TextDataRepository
	binaryDataRepo repository.BinaryDataRepository
	cardRepo       repository.CardRepository
	logger         *zap.Logger
}

func New(
	credentialRepo repository.CredentialRepository,
	textDataRepo repository.TextDataRepository,
	binaryDataRepo repository.BinaryDataRepository,
	cardRepo repository.CardRepository,
	logger *zap.Logger,
) SyncUsecase {
	return &syncUsecase{
		credentialRepo: credentialRepo,
		textDataRepo:   textDataRepo,
		binaryDataRepo: binaryDataRepo,
		cardRepo:       cardRepo,
		logger:         logger,
	}
}

func (u *syncUsecase) Execute(ctx context.Context, userID string) (*SyncResponse, error) {
	credentials, err := u.credentialRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get credentials", zap.Error(err))
		return nil, err
	}

	textDataList, err := u.textDataRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get text data", zap.Error(err))
		return nil, err
	}

	binaryDataList, err := u.binaryDataRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get binary data", zap.Error(err))
		return nil, err
	}

	cards, err := u.cardRepo.GetByUserID(ctx, userID)
	if err != nil {
		u.logger.Error("failed to get cards", zap.Error(err))
		return nil, err
	}

	now := time.Now()
	return &SyncResponse{
		Credentials: credentials,
		TextData:    textDataList,
		BinaryData:  binaryDataList,
		Cards:       cards,
		LastSync:    now,
	}, nil
}
