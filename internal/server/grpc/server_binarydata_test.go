package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	binarydataCreateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/create/mocks"
	binarydataDeleteMocks "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/delete/mocks"
	binarydataGetMocks "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/get/mocks"
	binarydataListMocks "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/list/mocks"
	binarydataUpdateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/update/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protobuf "google.golang.org/protobuf/proto"
)

func TestServer_CreateBinaryData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockCreate := binarydataCreateMocks.NewCreateBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, mockCreate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockCreate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.BinaryData")).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.BinaryDataRequest_builder{
		Title: protobuf.String("Test Binary"),
		Data:  []byte("test data"),
		Meta:  protobuf.String("test meta"),
	}.Build()

	resp, err := server.CreateBinaryData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_CreateBinaryData_BinaryDataTooLarge(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockCreate := binarydataCreateMocks.NewCreateBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, mockCreate, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockCreate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.BinaryData")).Return(model.ErrBinaryDataTooLarge).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.BinaryDataRequest_builder{
		Title: protobuf.String("Test Binary"),
		Data:  make([]byte, 11*1024*1024), // 11MB
		Meta:  protobuf.String("test meta"),
	}.Build()

	resp, err := server.CreateBinaryData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_GetBinaryData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := binarydataGetMocks.NewGetBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, mockGet, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	binaryDataID := "binary-456"

	now := time.Now()
	binaryData := &model.BinaryData{
		ID:            binaryDataID,
		UserID:        userID,
		Title:         "Test Binary",
		DataEncrypted: []byte("encrypted data"),
		Meta:          "test meta",
		Version:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockGet.EXPECT().Execute(mock.Anything, userID, binaryDataID).Return(binaryData, nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{
		Id: protobuf.String(binaryDataID),
	}.Build()

	resp, err := server.GetBinaryData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, binaryDataID, resp.GetId())
	assert.Equal(t, "Test Binary", resp.GetTitle())
}

func TestServer_UpdateBinaryData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := binarydataUpdateMocks.NewUpdateBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, mockUpdate, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	binaryDataID := "binary-456"

	mockUpdate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.BinaryData")).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateBinaryDataRequest_builder{
		Id:      protobuf.String(binaryDataID),
		Title:   protobuf.String("Updated"),
		Data:    []byte("updated data"),
		Meta:    protobuf.String("updated meta"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateBinaryData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_UpdateBinaryData_VersionConflict(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := binarydataUpdateMocks.NewUpdateBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, mockUpdate, nil, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	binaryDataID := "binary-456"

	mockUpdate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.BinaryData")).Return(model.ErrVersionConflict).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateBinaryDataRequest_builder{
		Id:      protobuf.String(binaryDataID),
		Title:   protobuf.String("Updated"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateBinaryData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Aborted, status.Code(err))
}

func TestServer_DeleteBinaryData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := binarydataDeleteMocks.NewDeleteBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, mockDelete, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	binaryDataID := "binary-456"

	mockDelete.EXPECT().Execute(mock.Anything, userID, binaryDataID).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(binaryDataID),
	}.Build()

	resp, err := server.DeleteBinaryData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DeleteBinaryData_NotFound(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := binarydataDeleteMocks.NewDeleteBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, mockDelete, nil, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	binaryDataID := "nonexistent"

	mockDelete.EXPECT().Execute(mock.Anything, userID, binaryDataID).Return(model.ErrBinaryDataNotFound).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(binaryDataID),
	}.Build()

	resp, err := server.DeleteBinaryData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestServer_ListBinaryData_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := binarydataListMocks.NewListBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, mockList, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	now := time.Now()
	binaryDataList := []*model.BinaryData{
		{
			ID:            "binary-1",
			UserID:        userID,
			Title:         "Binary 1",
			DataEncrypted: []byte("encrypted data 1"),
			Meta:          "meta 1",
			Version:       1,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	mockList.EXPECT().Execute(mock.Anything, userID).Return(binaryDataList, nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := &proto.ListBinaryDataRequest{}

	resp, err := server.ListBinaryData(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.GetBinaryData(), 1)
	assert.Equal(t, "binary-1", resp.GetBinaryData()[0].GetId())
}

func TestServer_ListBinaryData_InternalError(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := binarydataListMocks.NewListBinaryDataUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, mockList, nil, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockList.EXPECT().Execute(mock.Anything, userID).Return(nil, assert.AnError).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := &proto.ListBinaryDataRequest{}

	resp, err := server.ListBinaryData(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
}
