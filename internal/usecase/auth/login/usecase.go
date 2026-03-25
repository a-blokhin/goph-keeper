package login

import (
	"context"

	"github.com/a-blokhin/goph-keeper/internal/crypto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"go.uber.org/zap"
)

type loginUsecase struct {
	userRepo repository.UserRepository
	jwt      *jwt.JWTManager
	logger   *zap.Logger
}

func New(userRepo repository.UserRepository, jwtManager *jwt.JWTManager, logger *zap.Logger) LoginUsecase {
	return &loginUsecase{
		userRepo: userRepo,
		jwt:      jwtManager,
		logger:   logger,
	}
}

func (u *loginUsecase) Execute(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		u.logger.Error("failed to get user by email", zap.Error(err))
		return nil, "", model.ErrInvalidCredentials
	}

	if !crypto.VerifyPassword(password, user.PasswordHash, user.Salt) {
		return nil, "", model.ErrInvalidCredentials
	}

	token, err := u.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		u.logger.Error("failed to generate token", zap.Error(err))
		return nil, "", err
	}

	return user, token, nil
}
