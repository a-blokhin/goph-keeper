package login

import (
	"context"
	"errors"
	"testing"

	"github.com/a-blokhin/goph-keeper/internal/crypto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupUser(t *testing.T, password string) *model.User {
	t.Helper()
	hash, salt, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return &model.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: hash,
		Salt:         salt,
	}
}

func TestLoginUsecase_Execute_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	storedUser := setupUser(t, "password123")

	userRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").
		Return(storedUser, nil)

	user, token, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user-123", user.ID)
	assert.NotEmpty(t, token)

	claims, err := jwtManager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestLoginUsecase_Execute_WrongPassword(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	storedUser := setupUser(t, "password123")

	userRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").
		Return(storedUser, nil)

	user, token, err := uc.Execute(context.Background(), "test@example.com", "wrongpassword")

	assert.ErrorIs(t, err, model.ErrInvalidCredentials)
	assert.Nil(t, user)
	assert.Empty(t, token)
}

func TestLoginUsecase_Execute_UserNotFound(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	jwtManager := jwt.NewJWTManager("test-secret-key-12345678901234567890")
	logger := zap.NewNop()

	uc := New(userRepo, jwtManager, logger)

	userRepo.EXPECT().GetByEmail(mock.Anything, "nonexistent@example.com").
		Return(nil, errors.New("user not found"))

	user, token, err := uc.Execute(context.Background(), "nonexistent@example.com", "password123")

	assert.ErrorIs(t, err, model.ErrInvalidCredentials)
	assert.Nil(t, user)
	assert.Empty(t, token)
}
