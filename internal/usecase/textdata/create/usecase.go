package create

import (
	"context"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type createTextDataUsecase struct {
	textDataRepo      repository.TextDataRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(textDataRepo repository.TextDataRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) CreateTextDataUsecase {
	return &createTextDataUsecase{
		textDataRepo:      textDataRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *createTextDataUsecase) Execute(ctx context.Context, userID string, textData *model.TextData) error {
	if err := u.encryptionService.EncryptTextData(textData); err != nil {
		u.logger.Error("failed to encrypt text data", zap.Error(err))
		return err
	}

	textData.UserID = userID
	textData.Version = 1
	now := time.Now()
	textData.CreatedAt = now
	textData.UpdatedAt = now

	err := u.textDataRepo.Create(ctx, textData)
	if err != nil {
		u.logger.Error("failed to create text data", zap.Error(err))
		return err
	}

	return nil
}
