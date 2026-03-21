package binarydataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
)

type getBinaryDataUsecase struct {
	binaryDataRepo repository.BinaryDataRepository
}

func New(binaryDataRepo repository.BinaryDataRepository) GetBinaryDataUsecase {
	return &getBinaryDataUsecase{
		binaryDataRepo: binaryDataRepo,
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

	return binaryData, nil
}
