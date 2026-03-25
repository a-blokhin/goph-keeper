package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	cardCreateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/card/create/mocks"
	cardDeleteMocks "github.com/a-blokhin/goph-keeper/internal/usecase/card/delete/mocks"
	cardGetMocks "github.com/a-blokhin/goph-keeper/internal/usecase/card/get/mocks"
	cardListMocks "github.com/a-blokhin/goph-keeper/internal/usecase/card/list/mocks"
	cardUpdateMocks "github.com/a-blokhin/goph-keeper/internal/usecase/card/update/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protobuf "google.golang.org/protobuf/proto"
)

func TestServer_CreateCard_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockCreate := cardCreateMocks.NewCreateCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockCreate, nil, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockCreate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.Card")).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.CardRequest_builder{
		Title:      protobuf.String("Test Card"),
		CardNumber: protobuf.String("1234567890123456"),
		CardHolder: protobuf.String("John Doe"),
		Expiry:     protobuf.String("12/25"),
		Cvv:        protobuf.String("123"),
		Meta:       protobuf.String("test meta"),
	}.Build()

	resp, err := server.CreateCard(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_GetCard_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := cardGetMocks.NewGetCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockGet, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	cardID := "card-456"

	now := time.Now()
	card := &model.Card{
		ID:                  cardID,
		UserID:              userID,
		Title:               "Test Card",
		CardNumberEncrypted: "encrypted-number",
		CardHolderEncrypted: "encrypted-holder",
		ExpiryEncrypted:     "encrypted-expiry",
		CVVEncrypted:        "encrypted-cvv",
		Meta:                "test meta",
		Version:             1,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	mockGet.EXPECT().Execute(mock.Anything, userID, cardID).Return(card, nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{
		Id: protobuf.String(cardID),
	}.Build()

	resp, err := server.GetCard(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, cardID, resp.GetId())
	assert.Equal(t, "Test Card", resp.GetTitle())
}

func TestServer_GetCard_NotFound(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockGet := cardGetMocks.NewGetCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockGet, nil, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	cardID := "nonexistent"

	mockGet.EXPECT().Execute(mock.Anything, userID, cardID).Return(nil, model.ErrCardNotFound).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.GetRequest_builder{
		Id: protobuf.String(cardID),
	}.Build()

	resp, err := server.GetCard(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestServer_UpdateCard_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := cardUpdateMocks.NewUpdateCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockUpdate, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	cardID := "card-456"

	mockUpdate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.Card")).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateCardRequest_builder{
		Id:         protobuf.String(cardID),
		Title:      protobuf.String("Updated Card"),
		CardNumber: protobuf.String("9876543210987654"),
		CardHolder: protobuf.String("Jane Doe"),
		Expiry:     protobuf.String("06/26"),
		Cvv:        protobuf.String("456"),
		Meta:       protobuf.String("updated meta"),
		Version:    protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateCard(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_UpdateCard_VersionConflict(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockUpdate := cardUpdateMocks.NewUpdateCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockUpdate, nil, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	cardID := "card-456"

	mockUpdate.EXPECT().Execute(mock.Anything, userID, mock.AnythingOfType("*model.Card")).Return(model.ErrVersionConflict).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.UpdateCardRequest_builder{
		Id:      protobuf.String(cardID),
		Title:   protobuf.String("Updated Card"),
		Version: protobuf.Int32(1),
	}.Build()

	resp, err := server.UpdateCard(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Aborted, status.Code(err))
}

func TestServer_DeleteCard_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := cardDeleteMocks.NewDeleteCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockDelete, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	cardID := "card-456"

	mockDelete.EXPECT().Execute(mock.Anything, userID, cardID).Return(nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(cardID),
	}.Build()

	resp, err := server.DeleteCard(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServer_DeleteCard_NotFound(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockDelete := cardDeleteMocks.NewDeleteCardUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockDelete, nil, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")
	cardID := "nonexistent"

	mockDelete.EXPECT().Execute(mock.Anything, userID, cardID).Return(model.ErrCardNotFound).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(cardID),
	}.Build()

	resp, err := server.DeleteCard(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestServer_ListCards_Success(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := cardListMocks.NewListCardsUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockList, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	now := time.Now()
	cardList := []*model.Card{
		{
			ID:                  "card-1",
			UserID:              userID,
			Title:               "Card 1",
			CardNumberEncrypted: "encrypted-number-1",
			CardHolderEncrypted: "encrypted-holder-1",
			ExpiryEncrypted:     "encrypted-expiry-1",
			CVVEncrypted:        "encrypted-cvv-1",
			Meta:                "meta 1",
			Version:             1,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	mockList.EXPECT().Execute(mock.Anything, userID).Return(cardList, nil).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := &proto.ListCardsRequest{}

	resp, err := server.ListCards(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.GetCards(), 1)
	assert.Equal(t, "card-1", resp.GetCards()[0].GetId())
}

func TestServer_ListCards_InternalError(t *testing.T) {
	logger := zap.NewNop()
	jwtManager := jwt.NewJWTManager("test-secret-key")

	mockList := cardListMocks.NewListCardsUsecase(t)
	server := NewServer(logger, jwtManager, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, mockList, nil)

	userID := "user-123"
	token, _ := jwtManager.GenerateToken(userID, "test@example.com")

	mockList.EXPECT().Execute(mock.Anything, userID).Return(nil, assert.AnError).Once()

	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.New(map[string]string{"authorization": "Bearer " + token}))

	req := &proto.ListCardsRequest{}

	resp, err := server.ListCards(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
}
