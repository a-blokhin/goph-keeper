package update

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
	"go.uber.org/zap"
)

type updateTextDataUsecase struct {
	textDataRepo      repository.TextDataRepository
	encryptionService encryption.EncryptionService
	logger            *zap.Logger
}

func New(textDataRepo repository.TextDataRepository, encryptionService encryption.EncryptionService, logger *zap.Logger) UpdateTextDataUsecase {
	return &updateTextDataUsecase{
		textDataRepo:      textDataRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

func (u *updateTextDataUsecase) Execute(ctx context.Context, userID string, textData *model.TextData) error {
	existing, err := u.textDataRepo.GetByID(ctx, textData.ID)
	if err != nil {
		u.logger.Error("failed to get text data", zap.Error(err))
		return err
	}

	if existing.UserID != userID {
		return model.ErrForbidden
	}

	if existing.Version != textData.Version {
		return model.ErrVersionConflict
	}

	if err := u.encryptionService.EncryptTextData(textData); err != nil {
		u.logger.Error("failed to encrypt text data", zap.Error(err))
		return err
	}

	textData.UserID = userID
	textData.Version = existing.Version + 1

	err = u.textDataRepo.Update(ctx, textData)
	if err != nil {
		u.logger.Error("failed to update text data", zap.Error(err))
		return err
	}

	return nil
}
