package register

import (
	"context"
	"errors"

	"github.com/a-blokhin/goph-keeper/internal/crypto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type registerUsecase struct {
	userRepo repository.UserRepository
	jwt      *jwt.JWTManager
	logger   *zap.Logger
}

func New(userRepo repository.UserRepository, jwtManager *jwt.JWTManager, logger *zap.Logger) RegisterUsecase {
	return &registerUsecase{
		userRepo: userRepo,
		jwt:      jwtManager,
		logger:   logger,
	}
}

func (u *registerUsecase) Execute(ctx context.Context, email, password string) (*model.User, string, error) {
	passwordHash, salt, err := crypto.HashPassword(password)
	if err != nil {
		u.logger.Error("failed to hash password", zap.Error(err))
		return nil, "", err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: passwordHash,
		Salt:         salt,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, "", err
		}
		u.logger.Error("failed to create user", zap.Error(err))
		return nil, "", err
	}

	token, err := u.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		u.logger.Error("failed to generate token", zap.Error(err))
		return nil, "", err
	}

	return user, token, nil
}
