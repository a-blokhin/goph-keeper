package binarydataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/service/encryption"
)

type getBinaryDataUsecase struct {
	binaryDataRepo    repository.BinaryDataRepository
	encryptionService encryption.EncryptionService
}

func New(binaryDataRepo repository.BinaryDataRepository, encryptionService encryption.EncryptionService) GetBinaryDataUsecase {
	return &getBinaryDataUsecase{
		binaryDataRepo:    binaryDataRepo,
		encryptionService: encryptionService,
	}
}

func (u *getBinaryDataUsecase) Execute(ctx context.Context, userID, id string) (*model.BinaryData, error) {
	binaryData, err := u.binaryDataRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if binaryData.UserID != userID {
		return nil, model.ErrBinaryDataNotFound
	}

	if err := u.encryptionService.DecryptBinaryData(binaryData); err != nil {
		return nil, err
	}

	return binaryData, nil
}
