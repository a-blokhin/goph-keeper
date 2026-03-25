package grpc

import (
	"context"
	"testing"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	loginMocks "github.com/a-blokhin/goph-keeper/internal/usecase/auth/login/mocks"
	registerMocks "github.com/a-blokhin/goph-keeper/internal/usecase/auth/register/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	protobuf "google.golang.org/protobuf/proto"
)

func TestServer_Register_Success(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockRegister := registerMocks.NewRegisterUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		mockRegister, nil, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		nil, // sync usecase
	)

	// Test data
	email := "test@example.com"
	password := "password123"
	userID := "user-123"
	token := "test-token"

	// Setup expectations using EXPECT() pattern
	mockRegister.EXPECT().
		Execute(mock.Anything, email, password).
		Return(&model.User{ID: userID, Email: email}, token, nil).
		Once()

	// Execute
	ctx := context.Background()
	req := proto.RegisterRequest_builder{
		Email:    protobuf.String(email),
		Password: protobuf.String(password),
	}.Build()

	resp, err := server.Register(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, token, resp.GetToken())
}

func TestServer_Register_InternalError(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockRegister := registerMocks.NewRegisterUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		mockRegister, nil, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		nil, // sync usecase
	)

	// Test data
	email := "test@example.com"
	password := "password123"

	// Setup expectations using EXPECT() pattern
	mockRegister.EXPECT().
		Execute(mock.Anything, email, password).
		Return(nil, "", status.Error(codes.Internal, "internal error")).
		Once()

	// Execute
	ctx := context.Background()
	req := proto.RegisterRequest_builder{
		Email:    protobuf.String(email),
		Password: protobuf.String(password),
	}.Build()

	resp, err := server.Register(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
}

func TestServer_Login_Success(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockLogin := loginMocks.NewLoginUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		nil, mockLogin, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		nil, // sync usecase
	)

	// Test data
	email := "test@example.com"
	password := "password123"
	userID := "user-123"
	token := "test-token"

	// Setup expectations using EXPECT() pattern
	mockLogin.EXPECT().
		Execute(mock.Anything, email, password).
		Return(&model.User{ID: userID, Email: email}, token, nil).
		Once()

	// Execute
	ctx := context.Background()
	req := proto.LoginRequest_builder{
		Email:    protobuf.String(email),
		Password: protobuf.String(password),
	}.Build()

	resp, err := server.Login(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, token, resp.GetToken())
}

func TestServer_Login_InternalError(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockLogin := loginMocks.NewLoginUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		nil, mockLogin, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		nil, // sync usecase
	)

	// Test data
	email := "test@example.com"
	password := "password123"

	// Setup expectations using EXPECT() pattern
	mockLogin.EXPECT().
		Execute(mock.Anything, email, password).
		Return(nil, "", status.Error(codes.Internal, "internal error")).
		Once()

	// Execute
	ctx := context.Background()
	req := proto.LoginRequest_builder{
		Email:    protobuf.String(email),
		Password: protobuf.String(password),
	}.Build()

	resp, err := server.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
}
