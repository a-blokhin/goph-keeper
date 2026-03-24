package sync

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type syncUsecase struct {
	credentialRepo    repository.CredentialRepository
	textDataRepo      repository.TextDataRepository
	binaryDataRepo    repository.BinaryDataRepository
	cardRepo          repository.CardRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(
	credentialRepo repository.CredentialRepository,
	textDataRepo repository.TextDataRepository,
	binaryDataRepo repository.BinaryDataRepository,
	cardRepo repository.CardRepository,
	encryptionService encryption.EncryptionService,
	logger *zap.Logger,
) SyncUsecase {
	return &syncUsecase{
		credentialRepo:    credentialRepo,
		textDataRepo:      textDataRepo,
		binaryDataRepo:    binaryDataRepo,
		cardRepo:          cardRepo,
		encryptionService: encryptionService,
		logger:            logger,
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

	for _, cred := range credentials {
		if err := u.encryptionService.DecryptCredential(cred); err != nil {
			u.logger.Error("failed to decrypt credential", zap.Error(err))
			return nil, err
		}
	}

	for _, td := range textDataList {
		if err := u.encryptionService.DecryptTextData(td); err != nil {
			u.logger.Error("failed to decrypt text data", zap.Error(err))
			return nil, err
		}
	}

	for _, bd := range binaryDataList {
		if err := u.encryptionService.DecryptBinaryData(bd); err != nil {
			u.logger.Error("failed to decrypt binary data", zap.Error(err))
			return nil, err
		}
	}

	for _, card := range cards {
		if err := u.encryptionService.DecryptCardData(card); err != nil {
			u.logger.Error("failed to decrypt card data", zap.Error(err))
			return nil, err
		}
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
