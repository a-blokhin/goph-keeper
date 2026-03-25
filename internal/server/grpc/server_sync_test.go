package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/usecase/sync"
	syncMocks "github.com/a-blokhin/goph-keeper/internal/usecase/sync/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestServer_Sync_Success(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockSync := syncMocks.NewSyncUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		nil, nil, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		mockSync, // sync usecase
	)

	// Test data
	userID := "user-123"
	token, err := jwtManager.GenerateToken(userID, "test@example.com")
	assert.NoError(t, err)

	now := time.Now()
	syncResponse := &sync.SyncResponse{
		Credentials: []*model.Credential{
			{
				ID:                "cred-1",
				UserID:            userID,
				Title:             "Test Credential",
				Login:             "testuser",
				PasswordEncrypted: "testpass",
				Meta:              "test meta",
				Version:           1,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		TextData: []*model.TextData{
			{
				ID:            "text-1",
				UserID:        userID,
				Title:         "Test Text",
				DataEncrypted: "test data",
				Meta:          "test meta",
				Version:       1,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
		BinaryData: []*model.BinaryData{
			{
				ID:            "binary-1",
				UserID:        userID,
				Title:         "Test Binary",
				DataEncrypted: []byte("encrypted data"),
				Meta:          "test meta",
				Version:       1,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
		Cards: []*model.Card{
			{
				ID:                  "card-1",
				UserID:              userID,
				Title:               "Test Card",
				CardNumberEncrypted: "1234567890123456",
				CardHolderEncrypted: "Test Holder",
				CVVEncrypted:        "123",
				ExpiryEncrypted:     "12/25",
				Meta:                "test meta",
				Version:             1,
				CreatedAt:           now,
				UpdatedAt:           now,
			},
		},
	}

	// Setup expectations using EXPECT() pattern
	mockSync.EXPECT().
		Execute(mock.Anything, userID).
		Return(syncResponse, nil).
		Once()

	// Execute
	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewIncomingContext(ctx, md)

	req := proto.SyncRequest_builder{}.Build()

	resp, err := server.Sync(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.GetCredentials(), 1)
	assert.Len(t, resp.GetTextData(), 1)
	assert.Len(t, resp.GetBinaryData(), 1)
	assert.Len(t, resp.GetCards(), 1)
}

func TestServer_Sync_Unauthenticated(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockSync := syncMocks.NewSyncUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		nil, nil, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		mockSync, // sync usecase
	)

	// Execute without authentication
	ctx := context.Background()
	req := proto.SyncRequest_builder{}.Build()

	resp, err := server.Sync(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, statusErr.Code())
}

func TestServer_Sync_InternalError(t *testing.T) {
	// Setup
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockSync := syncMocks.NewSyncUsecase(t)
	server := NewServer(
		logger,
		jwtManager,
		nil, nil, // auth usecases
		nil, nil, nil, nil, nil, // credential usecases
		nil, nil, nil, nil, nil, // text data usecases
		nil, nil, nil, nil, nil, // binary data usecases
		nil, nil, nil, nil, nil, // card usecases
		mockSync, // sync usecase
	)

	// Test data
	userID := "user-123"
	token, err := jwtManager.GenerateToken(userID, "test@example.com")
	assert.NoError(t, err)

	// Setup expectations using EXPECT() pattern
	mockSync.EXPECT().
		Execute(mock.Anything, userID).
		Return(nil, assert.AnError).
		Once()

	// Execute
	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewIncomingContext(ctx, md)

	req := proto.SyncRequest_builder{}.Build()

	resp, err := server.Sync(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
}
