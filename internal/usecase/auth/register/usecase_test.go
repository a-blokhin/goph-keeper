package register

import (
	"context"
	"errors"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestRegisterUsecase_Execute_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).
		Run(func(ctx context.Context, user *model.User) {
			user.ID = "generated-id"
		}).
		Return(nil)

	user, token, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEmpty(t, user.Salt)

	claims, err := jwtManager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "generated-id", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestRegisterUsecase_Execute_UserAlreadyExists(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).
		Return(model.ErrUserAlreadyExists)

	user, token, err := uc.Execute(context.Background(), "existing@example.com", "password123")

	assert.ErrorIs(t, err, model.ErrUserAlreadyExists)
	assert.Nil(t, user)
	assert.Empty(t, token)
}

func TestRegisterUsecase_Execute_RepoError(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	repoErr := errors.New("database connection failed")
	userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).
		Return(repoErr)

	user, token, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, repoErr, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
}

func TestRegisterUsecase_Execute_PasswordIsHashed(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	var capturedUser *model.User
	userRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.User")).
		Run(func(ctx context.Context, user *model.User) {
			capturedUser = user
		}).
		Return(nil)

	_, _, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEqual(t, "password123", capturedUser.PasswordHash)
	assert.NotEmpty(t, capturedUser.Salt)
}
