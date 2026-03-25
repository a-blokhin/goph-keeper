package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	credentialCreateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/credential/create/mocks"
	credentialDeleteMocks "github.com/a-blokhin/goph-keeper/internal/usecase/credential/delete/mocks"
	credentialGetMocks "github.com/a-blokhin/goph-keeper/internal/usecase/credential/get/mocks"
	credentialListMocks "github.com/a-blokhin/goph-keeper/internal/usecase/credential/list/mocks"
	credentialUpdateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/credential/update/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protobuf "google.golang.org/protobuf/proto"
)

func TestServer_CreateCredential_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockCreate := credentialCreateMocks.NewCreateCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		mockCreate, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockCreate.EXPECT().
		Execute(mock.Anything, userID, mock.AnythingOfType("*model.Credential")).
		Return(nil).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.CredentialRequest_builder{
		Title:    protobuf.String("Test Credential"),
		Login:    protobuf.String("testuser"),
		Password: protobuf.String("testpass"),
	}.Build()

	resp, err := server.CreateCredential(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_GetCredential_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := credentialGetMocks.NewGetCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, mockGet, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	credID := "cred-123"

	credential := &model.Credential{
		ID:                credID,
		UserID:            userID,
		Title:             "Test",
		Login:             "testuser",
		PasswordEncrypted: "encrypted",
		Version:           1,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	mockGet.EXPECT().
		Execute(mock.Anything, userID, credID).
		Return(credential, nil).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{Id: protobuf.String(credID)}.Build()

	resp, err := server.GetCredential(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, credID, resp.GetId())
}

func TestServer_GetCredential_NotFound(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := credentialGetMocks.NewGetCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, mockGet, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	credID := "cred-123"

	mockGet.EXPECT().
		Execute(mock.Anything, userID, credID).
		Return(nil, model.ErrCredentialNotFound).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{Id: protobuf.String(credID)}.Build()

	resp, err := server.GetCredential(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())
}

func TestServer_UpdateCredential_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := credentialUpdateMocks.NewUpdateCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, mockUpdate, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockUpdate.EXPECT().
		Execute(mock.Anything, userID, mock.AnythingOfType("*model.Credential")).
		Return(nil).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateCredentialRequest_builder{
		Id:      protobuf.String("cred-123"),
		Title:   protobuf.String("Updated"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateCredential(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_UpdateCredential_VersionConflict(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := credentialUpdateMocks.NewUpdateCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, mockUpdate, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockUpdate.EXPECT().
		Execute(mock.Anything, userID, mock.AnythingOfType("*model.Credential")).
		Return(model.ErrVersionConflict).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateCredentialRequest_builder{
		Id:      protobuf.String("cred-123"),
		Title:   protobuf.String("Updated"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateCredential(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Aborted, statusErr.Code())
}

func TestServer_DeleteCredential_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := credentialDeleteMocks.NewDeleteCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, nil, mockDelete, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	credID := "cred-123"

	mockDelete.EXPECT().
		Execute(mock.Anything, userID, credID).
		Return(nil).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{Id: protobuf.String(credID)}.Build()

	resp, err := server.DeleteCredential(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DeleteCredential_NotFound(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := credentialDeleteMocks.NewDeleteCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, nil, mockDelete, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	credID := "cred-123"

	mockDelete.EXPECT().
		Execute(mock.Anything, userID, credID).
		Return(model.ErrCredentialNotFound).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{Id: protobuf.String(credID)}.Build()

	resp, err := server.DeleteCredential(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())
}

func TestServer_DeleteCredential_Forbidden(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := credentialDeleteMocks.NewDeleteCredentialUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, nil, mockDelete, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	credID := "cred-123"

	mockDelete.EXPECT().
		Execute(mock.Anything, userID, credID).
		Return(model.ErrForbidden).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{Id: protobuf.String(credID)}.Build()

	resp, err := server.DeleteCredential(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, statusErr.Code())
}

func TestServer_ListCredentials_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := credentialListMocks.NewListCredentialsUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, nil, nil, mockList,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	credentials := []*model.Credential{
		{
			ID:                "cred-1",
			UserID:            userID,
			Title:             "Test 1",
			Login:             "user1",
			PasswordEncrypted: "encrypted1",
			Version:           1,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                "cred-2",
			UserID:            userID,
			Title:             "Test 2",
			Login:             "user2",
			PasswordEncrypted: "encrypted2",
			Version:           1,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}

	mockList.EXPECT().
		Execute(mock.Anything, userID).
		Return(credentials, nil).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.ListCredentialsRequest_builder{}.Build()

	resp, err := server.ListCredentials(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.GetCredentials(), 2)
}

func TestServer_ListCredentials_InternalError(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := credentialListMocks.NewListCredentialsUsecase(t)
	server := NewServer(
		logger, jwtManager,
		nil, nil,
		nil, nil, nil, nil, mockList,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil,
		nil,
	)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockList.EXPECT().
		Execute(mock.Anything, userID).
		Return(nil, errors.New("internal error")).
		Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.ListCredentialsRequest_builder{}.Build()

	resp, err := server.ListCredentials(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
}
