package textdataget

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
)

type getTextDataUsecase struct {
	textDataRepo repository.TextDataRepository
}

func New(textDataRepo repository.TextDataRepository) GetTextDataUsecase {
	return &getTextDataUsecase{
		textDataRepo: textDataRepo,
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

	return textData, nil
}
