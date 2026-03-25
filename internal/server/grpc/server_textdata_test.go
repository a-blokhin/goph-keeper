package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	textdataCreateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/create/mocks"
	textdataDeleteMocks "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/delete/mocks"
	textdataGetMocks "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/get/mocks"
	textdataListMocks "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/list/mocks"
	textdataUpdateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/update/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protobuf "google.golang.org/protobuf/proto"
)

func TestServer_CreateTextData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockCreate := textdataCreateMocks.NewCreateTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		mockCreate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockCreate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.TextData")).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.TextDataRequest_builder{
		Title: protobuf.String("Test Text"),
		Data:  protobuf.String("test data"),
	}.Build()

	resp, err := server.CreateTextData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_GetTextData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := textdataGetMocks.NewGetTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, mockGet, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	textID := "text-123"

	textData := &model.TextData{
		ID:            textID,
		UserID:        userID,
		Title:         "Test",
		DataEncrypted: "encrypted",
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mockGet.EXPECT().Execute(mock.Anything, userID, textID).Return(textData, nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{Id: protobuf.String(textID)}.Build()

	resp, err := server.GetTextData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, textID, resp.GetId())
}

func TestServer_GetTextData_NotFound(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := textdataGetMocks.NewGetTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, mockGet, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	textID := "text-123"

	mockGet.EXPECT().Execute(mock.Anything, userID, textID).Return(nil, model.ErrTextDataNotFound).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{Id: protobuf.String(textID)}.Build()

	resp, err := server.GetTextData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())
}

func TestServer_UpdateTextData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := textdataUpdateMocks.NewUpdateTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, mockUpdate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockUpdate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.TextData")).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateTextDataRequest_builder{
		Id:      protobuf.String("text-123"),
		Title:   protobuf.String("Updated"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateTextData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_UpdateTextData_VersionConflict(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := textdataUpdateMocks.NewUpdateTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, mockUpdate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockUpdate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.TextData")).Return(model.ErrVersionConflict).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateTextDataRequest_builder{
		Id:      protobuf.String("text-123"),
		Title:   protobuf.String("Updated"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateTextData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Aborted, statusErr.Code())
}

func TestServer_DeleteTextData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := textdataDeleteMocks.NewDeleteTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, mockDelete, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	textID := "text-123"

	mockDelete.EXPECT().Execute(mock.Anything, userID, textID).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{Id: protobuf.String(textID)}.Build()

	resp, err := server.DeleteTextData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DeleteTextData_Forbidden(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := textdataDeleteMocks.NewDeleteTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, mockDelete, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	textID := "text-123"

	mockDelete.EXPECT().Execute(mock.Anything, userID, textID).Return(model.ErrForbidden).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{Id: protobuf.String(textID)}.Build()

	resp, err := server.DeleteTextData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, statusErr.Code())
}

func TestServer_ListTextData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := textdataListMocks.NewListTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, mockList, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	textDataList := []*model.TextData{
		{ID: "text-1", UserID: userID, Title: "Test 1", DataEncrypted: "encrypted1", Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "text-2", UserID: userID, Title: "Test 2", DataEncrypted: "encrypted2", Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockList.EXPECT().Execute(mock.Anything, userID).Return(textDataList, nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.ListTextDataRequest_builder{}.Build()

	resp, err := server.ListTextData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.GetTextData(), 2)
}

func TestServer_ListTextData_InternalError(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := textdataListMocks.NewListTextDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, mockList, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockList.EXPECT().Execute(mock.Anything, userID).Return(nil, errors.New("internal error")).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.ListTextDataRequest_builder{}.Build()

	resp, err := server.ListTextData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())
}
