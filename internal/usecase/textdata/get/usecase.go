package textdataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
)

type getTextDataUsecase struct {
	textDataRepo      repository.TextDataRepository
	encryptionService encryption.EncryptionService
}

func New(textDataRepo repository.TextDataRepository, encryptionService encryption.EncryptionService) GetTextDataUsecase {
	return &getTextDataUsecase{
		textDataRepo:      textDataRepo,
		encryptionService: encryptionService,
	}
}

func (u *getTextDataUsecase) Execute(ctx context.Context, userID, id string) (*model.TextData, error) {
	textData, err := u.textDataRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if textData.UserID != userID {
		return nil, model.ErrTextDataNotFound
	}

	if err := u.encryptionService.DecryptTextData(textData); err != nil {
		return nil, err
	}

	return textData, nil
}
