//go:generate mockery --name=UserRepository --output=./mocks --outpkg=mocks --filename=user_repository_mock.go --with-expecter
//go:generate mockery --name=CredentialRepository --output=./mocks --outpkg=mocks --filename=credential_repository_mock.go --with-expecter
//go:generate mockery --name=TextDataRepository --output=./mocks --outpkg=mocks --filename=textdata_repository_mock.go --with-expecter
//go:generate mockery --name=BinaryDataRepository --output=./mocks --outpkg=mocks --filename=binarydata_repository_mock.go --with-expecter
//go:generate mockery --name=CardRepository --output=./mocks --outpkg=mocks --filename=card_repository_mock.go --with-expecter

package repository

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type CredentialRepository interface {
	Create(ctx context.Context, credential *model.Credential) error
	GetByID(ctx context.Context, id string) (*model.Credential, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Credential, error)
	Update(ctx context.Context, credential *model.Credential) error
	Delete(ctx context.Context, id string) error
}

type TextDataRepository interface {
	Create(ctx context.Context, textData *model.TextData) error
	GetByID(ctx context.Context, id string) (*model.TextData, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.TextData, error)
	Update(ctx context.Context, textData *model.TextData) error
	Delete(ctx context.Context, id string) error
}

type BinaryDataRepository interface {
	Create(ctx context.Context, binaryData *model.BinaryData) error
	GetByID(ctx context.Context, id string) (*model.BinaryData, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.BinaryData, error)
	Update(ctx context.Context, binaryData *model.BinaryData) error
	Delete(ctx context.Context, id string) error
}

type CardRepository interface {
	Create(ctx context.Context, card *model.Card) error
	GetByID(ctx context.Context, id string) (*model.Card, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Card, error)
	Update(ctx context.Context, card *model.Card) error
	Delete(ctx context.Context, id string) error
}
