package client

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/api/proto/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewClient(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.localStorage)
}

func TestClient_SaveToken(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Test saving a token
	token := "test-token-123"
	err = client.SaveToken(token)
	assert.NoError(t, err)

	// Verify token was saved
	assert.Equal(t, token, client.token)

	// Clean up the token file
	homeDir, _ := os.UserHomeDir()
	tokenPath := filepath.Join(homeDir, tokenFile)
	_ = os.Remove(tokenPath)
}

func TestClient_loadToken(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Test loading when token is already set
	client.token = "existing-token"
	err = client.loadToken()
	assert.NoError(t, err)
	assert.Equal(t, "existing-token", client.token)

	// Test loading from file
	token := "file-token-456"
	err = client.SaveToken(token)
	require.NoError(t, err)

	// Clear token and reload
	client.token = ""
	err = client.loadToken()
	assert.NoError(t, err)
	assert.Equal(t, token, client.token)

	// Clean up the token file
	homeDir, _ := os.UserHomeDir()
	tokenPath := filepath.Join(homeDir, tokenFile)
	_ = os.Remove(tokenPath)
}

func TestClient_Register(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	expectedToken := "registered-token-789"
	resp := proto.RegisterResponse_builder{
		Token: &expectedToken,
	}.Build()

	mockClient.EXPECT().Register(mock.Anything, mock.AnythingOfType("*proto.RegisterRequest"), mock.Anything).Return(resp, nil)

	// Test registration
	token, err := client.Register(context.Background(), "test@example.com", "password123")
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_Login(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	expectedToken := "login-token-789"
	resp := proto.LoginResponse_builder{
		Token: &expectedToken,
	}.Build()

	mockClient.EXPECT().Login(mock.Anything, mock.AnythingOfType("*proto.LoginRequest"), mock.Anything).Return(resp, nil)

	// Test login
	token, err := client.Login(context.Background(), "test@example.com", "password123")
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_CreateCredential(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.CredentialResponse_builder{
		Id:       stringPtr("cred-1"),
		Title:    stringPtr("Test Credential"),
		Login:    stringPtr("testuser"),
		Password: stringPtr("testpass"),
		Meta:     stringPtr("test meta"),
		Version:  int32Ptr(1),
	}.Build()

	mockClient.EXPECT().CreateCredential(mock.Anything, mock.AnythingOfType("*proto.CredentialRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test creating credential
	result, err := client.CreateCredential(context.Background(), "Test Credential", "testuser", "testpass", "test meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cred-1", result.GetId())
	assert.Equal(t, "Test Credential", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_GetCredential(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.CredentialResponse_builder{
		Id:       stringPtr("cred-1"),
		Title:    stringPtr("Test Credential"),
		Login:    stringPtr("testuser"),
		Password: stringPtr("testpass"),
		Meta:     stringPtr("test meta"),
		Version:  int32Ptr(1),
	}.Build()

	mockClient.On("GetCredential", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test getting credential
	result, err := client.GetCredential(context.Background(), "cred-1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cred-1", result.GetId())
	assert.Equal(t, "Test Credential", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_ListCredentials(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	creds := []*proto.CredentialResponse{
		proto.CredentialResponse_builder{
			Id:       stringPtr("cred-1"),
			Title:    stringPtr("Test Credential 1"),
			Login:    stringPtr("testuser1"),
			Password: stringPtr("testpass1"),
			Meta:     stringPtr("test meta 1"),
			Version:  int32Ptr(1),
		}.Build(),
		proto.CredentialResponse_builder{
			Id:       stringPtr("cred-2"),
			Title:    stringPtr("Test Credential 2"),
			Login:    stringPtr("testuser2"),
			Password: stringPtr("testpass2"),
			Meta:     stringPtr("test meta 2"),
			Version:  int32Ptr(1),
		}.Build(),
	}
	resp := proto.CredentialsListResponse_builder{
		Credentials: creds,
	}.Build()

	mockClient.On("ListCredentials", mock.Anything, mock.AnythingOfType("*proto.ListCredentialsRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test listing credentials
	result, err := client.ListCredentials(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "cred-1", result[0].GetId())
	assert.Equal(t, "cred-2", result[1].GetId())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_UpdateCredential(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectations
	// First mock for GetCredential call inside UpdateCredential
	getResp := proto.CredentialResponse_builder{
		Id:       stringPtr("cred-1"),
		Title:    stringPtr("Test Credential"),
		Login:    stringPtr("testuser"),
		Password: stringPtr("testpass"),
		Meta:     stringPtr("test meta"),
		Version:  int32Ptr(1),
	}.Build()
	mockClient.On("GetCredential", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(getResp, nil).Once()

	// Second mock for UpdateCredential call
	updateResp := proto.CredentialResponse_builder{
		Id:       stringPtr("cred-1"),
		Title:    stringPtr("Updated Credential"),
		Login:    stringPtr("updateduser"),
		Password: stringPtr("updatedpass"),
		Meta:     stringPtr("updated meta"),
		Version:  int32Ptr(2),
	}.Build()
	mockClient.On("UpdateCredential", mock.Anything, mock.AnythingOfType("*proto.UpdateCredentialRequest"), mock.Anything).Return(updateResp, nil).Once()

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test updating credential
	result, err := client.UpdateCredential(context.Background(), "cred-1", "Updated Credential", "updateduser", "updatedpass", "updated meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cred-1", result.GetId())
	assert.Equal(t, "Updated Credential", result.GetTitle())
	assert.Equal(t, int32(2), result.GetVersion())

	// Verify mocks were called
	mockClient.AssertExpectations(t)
}

func TestClient_DeleteCredential(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.DeleteResponse_builder{}.Build()
	mockClient.On("DeleteCredential", mock.Anything, mock.AnythingOfType("*proto.DeleteRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test deleting credential
	err = client.DeleteCredential(context.Background(), "cred-1")
	assert.NoError(t, err)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_CreateTextData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.TextDataResponse_builder{
		Id:      stringPtr("text-1"),
		Title:   stringPtr("Test Text"),
		Data:    stringPtr("test data content"),
		Meta:    stringPtr("test meta"),
		Version: int32Ptr(1),
	}.Build()

	mockClient.On("CreateTextData", mock.Anything, mock.AnythingOfType("*proto.TextDataRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test creating text data
	result, err := client.CreateTextData(context.Background(), "Test Text", "test data content", "test meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "text-1", result.GetId())
	assert.Equal(t, "Test Text", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_GetTextData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.TextDataResponse_builder{
		Id:      stringPtr("text-1"),
		Title:   stringPtr("Test Text"),
		Data:    stringPtr("test data content"),
		Meta:    stringPtr("test meta"),
		Version: int32Ptr(1),
	}.Build()

	mockClient.On("GetTextData", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test getting text data
	result, err := client.GetTextData(context.Background(), "text-1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "text-1", result.GetId())
	assert.Equal(t, "Test Text", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_ListTextData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	textData := []*proto.TextDataResponse{
		proto.TextDataResponse_builder{
			Id:      stringPtr("text-1"),
			Title:   stringPtr("Test Text 1"),
			Data:    stringPtr("test data content 1"),
			Meta:    stringPtr("test meta 1"),
			Version: int32Ptr(1),
		}.Build(),
		proto.TextDataResponse_builder{
			Id:      stringPtr("text-2"),
			Title:   stringPtr("Test Text 2"),
			Data:    stringPtr("test data content 2"),
			Meta:    stringPtr("test meta 2"),
			Version: int32Ptr(1),
		}.Build(),
	}
	resp := proto.TextDataListResponse_builder{
		TextData: textData,
	}.Build()

	mockClient.On("ListTextData", mock.Anything, mock.AnythingOfType("*proto.ListTextDataRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test listing text data
	result, err := client.ListTextData(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "text-1", result[0].GetId())
	assert.Equal(t, "text-2", result[1].GetId())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_UpdateTextData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectations
	// First mock for GetTextData call inside UpdateTextData
	getResp := proto.TextDataResponse_builder{
		Id:      stringPtr("text-1"),
		Title:   stringPtr("Test Text"),
		Data:    stringPtr("test data content"),
		Meta:    stringPtr("test meta"),
		Version: int32Ptr(1),
	}.Build()
	mockClient.On("GetTextData", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(getResp, nil).Once()

	// Second mock for UpdateTextData call
	updateResp := proto.TextDataResponse_builder{
		Id:      stringPtr("text-1"),
		Title:   stringPtr("Updated Text"),
		Data:    stringPtr("updated data content"),
		Meta:    stringPtr("updated meta"),
		Version: int32Ptr(2),
	}.Build()
	mockClient.On("UpdateTextData", mock.Anything, mock.AnythingOfType("*proto.UpdateTextDataRequest"), mock.Anything).Return(updateResp, nil).Once()

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test updating text data
	result, err := client.UpdateTextData(context.Background(), "text-1", "Updated Text", "updated data content", "updated meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "text-1", result.GetId())
	assert.Equal(t, "Updated Text", result.GetTitle())
	assert.Equal(t, int32(2), result.GetVersion())

	// Verify mocks were called
	mockClient.AssertExpectations(t)
}

func TestClient_DeleteTextData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.DeleteResponse_builder{}.Build()
	mockClient.On("DeleteTextData", mock.Anything, mock.AnythingOfType("*proto.DeleteRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test deleting text data
	err = client.DeleteTextData(context.Background(), "text-1")
	assert.NoError(t, err)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_CreateBinaryData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.BinaryDataResponse_builder{
		Id:      stringPtr("binary-1"),
		Title:   stringPtr("Test Binary"),
		Data:    []byte("test binary content"),
		Meta:    stringPtr("test meta"),
		Version: int32Ptr(1),
	}.Build()

	mockClient.On("CreateBinaryData", mock.Anything, mock.AnythingOfType("*proto.BinaryDataRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test creating binary data
	result, err := client.CreateBinaryData(context.Background(), "Test Binary", []byte("test binary content"), "test meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "binary-1", result.GetId())
	assert.Equal(t, "Test Binary", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_CreateBinaryData_TooLarge(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test creating binary data that's too large (over 10MB)
	largeData := make([]byte, 11*1024*1024) // 11MB
	_, err = client.CreateBinaryData(context.Background(), "Large Binary", largeData, "test meta")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data too large")
}

func TestClient_GetBinaryData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.BinaryDataResponse_builder{
		Id:      stringPtr("binary-1"),
		Title:   stringPtr("Test Binary"),
		Data:    []byte("test binary content"),
		Meta:    stringPtr("test meta"),
		Version: int32Ptr(1),
	}.Build()

	mockClient.On("GetBinaryData", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test getting binary data
	result, err := client.GetBinaryData(context.Background(), "binary-1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "binary-1", result.GetId())
	assert.Equal(t, "Test Binary", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_ListBinaryData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	binaryData := []*proto.BinaryDataResponse{
		proto.BinaryDataResponse_builder{
			Id:      stringPtr("binary-1"),
			Title:   stringPtr("Test Binary 1"),
			Data:    []byte("test binary content 1"),
			Meta:    stringPtr("test meta 1"),
			Version: int32Ptr(1),
		}.Build(),
		proto.BinaryDataResponse_builder{
			Id:      stringPtr("binary-2"),
			Title:   stringPtr("Test Binary 2"),
			Data:    []byte("test binary content 2"),
			Meta:    stringPtr("test meta 2"),
			Version: int32Ptr(1),
		}.Build(),
	}
	resp := proto.BinaryDataListResponse_builder{
		BinaryData: binaryData,
	}.Build()

	mockClient.On("ListBinaryData", mock.Anything, mock.AnythingOfType("*proto.ListBinaryDataRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test listing binary data
	result, err := client.ListBinaryData(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "binary-1", result[0].GetId())
	assert.Equal(t, "binary-2", result[1].GetId())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_UpdateBinaryData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectations
	// First mock for GetBinaryData call inside UpdateBinaryData
	getResp := proto.BinaryDataResponse_builder{
		Id:      stringPtr("binary-1"),
		Title:   stringPtr("Test Binary"),
		Data:    []byte("test binary content"),
		Meta:    stringPtr("test meta"),
		Version: int32Ptr(1),
	}.Build()
	mockClient.On("GetBinaryData", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(getResp, nil).Once()

	// Second mock for UpdateBinaryData call
	updateResp := proto.BinaryDataResponse_builder{
		Id:      stringPtr("binary-1"),
		Title:   stringPtr("Updated Binary"),
		Data:    []byte("updated binary content"),
		Meta:    stringPtr("updated meta"),
		Version: int32Ptr(2),
	}.Build()
	mockClient.On("UpdateBinaryData", mock.Anything, mock.AnythingOfType("*proto.UpdateBinaryDataRequest"), mock.Anything).Return(updateResp, nil).Once()

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test updating binary data
	result, err := client.UpdateBinaryData(context.Background(), "binary-1", "Updated Binary", []byte("updated binary content"), "updated meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "binary-1", result.GetId())
	assert.Equal(t, "Updated Binary", result.GetTitle())
	assert.Equal(t, int32(2), result.GetVersion())

	// Verify mocks were called
	mockClient.AssertExpectations(t)
}

func TestClient_UpdateBinaryData_TooLarge(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test updating binary data that's too large (over 10MB)
	largeData := make([]byte, 11*1024*1024) // 11MB
	_, err = client.UpdateBinaryData(context.Background(), "binary-1", "Large Binary", largeData, "test meta")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data too large")
}

func TestClient_DeleteBinaryData(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.DeleteResponse_builder{}.Build()
	mockClient.On("DeleteBinaryData", mock.Anything, mock.AnythingOfType("*proto.DeleteRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test deleting binary data
	err = client.DeleteBinaryData(context.Background(), "binary-1")
	assert.NoError(t, err)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_CreateCard(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.CardResponse_builder{
		Id:         stringPtr("card-1"),
		Title:      stringPtr("Test Card"),
		CardNumber: stringPtr("1234567890123456"),
		CardHolder: stringPtr("Test User"),
		Expiry:     stringPtr("12/25"),
		Cvv:        stringPtr("123"),
		Meta:       stringPtr("test meta"),
		Version:    int32Ptr(1),
	}.Build()

	mockClient.On("CreateCard", mock.Anything, mock.AnythingOfType("*proto.CardRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test creating card
	result, err := client.CreateCard(context.Background(), "Test Card", "1234567890123456", "Test User", "12/25", "123", "test meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "card-1", result.GetId())
	assert.Equal(t, "Test Card", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_GetCard(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.CardResponse_builder{
		Id:         stringPtr("card-1"),
		Title:      stringPtr("Test Card"),
		CardNumber: stringPtr("1234567890123456"),
		CardHolder: stringPtr("Test User"),
		Expiry:     stringPtr("12/25"),
		Cvv:        stringPtr("123"),
		Meta:       stringPtr("test meta"),
		Version:    int32Ptr(1),
	}.Build()

	mockClient.On("GetCard", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test getting card
	result, err := client.GetCard(context.Background(), "card-1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "card-1", result.GetId())
	assert.Equal(t, "Test Card", result.GetTitle())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_ListCards(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	cards := []*proto.CardResponse{
		proto.CardResponse_builder{
			Id:         stringPtr("card-1"),
			Title:      stringPtr("Test Card 1"),
			CardNumber: stringPtr("1234567890123456"),
			CardHolder: stringPtr("Test User 1"),
			Expiry:     stringPtr("12/25"),
			Cvv:        stringPtr("123"),
			Meta:       stringPtr("test meta 1"),
			Version:    int32Ptr(1),
		}.Build(),
		proto.CardResponse_builder{
			Id:         stringPtr("card-2"),
			Title:      stringPtr("Test Card 2"),
			CardNumber: stringPtr("6543210987654321"),
			CardHolder: stringPtr("Test User 2"),
			Expiry:     stringPtr("11/26"),
			Cvv:        stringPtr("456"),
			Meta:       stringPtr("test meta 2"),
			Version:    int32Ptr(1),
		}.Build(),
	}
	resp := proto.CardsListResponse_builder{
		Cards: cards,
	}.Build()

	mockClient.On("ListCards", mock.Anything, mock.AnythingOfType("*proto.ListCardsRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test listing cards
	result, err := client.ListCards(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "card-1", result[0].GetId())
	assert.Equal(t, "card-2", result[1].GetId())

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_UpdateCard(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectations
	// First mock for GetCard call inside UpdateCard
	getResp := proto.CardResponse_builder{
		Id:         stringPtr("card-1"),
		Title:      stringPtr("Test Card"),
		CardNumber: stringPtr("1234567890123456"),
		CardHolder: stringPtr("Test User"),
		Expiry:     stringPtr("12/25"),
		Cvv:        stringPtr("123"),
		Meta:       stringPtr("test meta"),
		Version:    int32Ptr(1),
	}.Build()
	mockClient.On("GetCard", mock.Anything, mock.AnythingOfType("*proto.GetRequest"), mock.Anything).Return(getResp, nil).Once()

	// Second mock for UpdateCard call
	updateResp := proto.CardResponse_builder{
		Id:         stringPtr("card-1"),
		Title:      stringPtr("Updated Card"),
		CardNumber: stringPtr("6543210987654321"),
		CardHolder: stringPtr("Updated User"),
		Expiry:     stringPtr("11/26"),
		Cvv:        stringPtr("456"),
		Meta:       stringPtr("updated meta"),
		Version:    int32Ptr(2),
	}.Build()
	mockClient.On("UpdateCard", mock.Anything, mock.AnythingOfType("*proto.UpdateCardRequest"), mock.Anything).Return(updateResp, nil).Once()

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test updating card
	result, err := client.UpdateCard(context.Background(), "card-1", "Updated Card", "6543210987654321", "Updated User", "11/26", "456", "updated meta")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "card-1", result.GetId())
	assert.Equal(t, "Updated Card", result.GetTitle())
	assert.Equal(t, int32(2), result.GetVersion())

	// Verify mocks were called
	mockClient.AssertExpectations(t)
}

func TestClient_DeleteCard(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	resp := proto.DeleteResponse_builder{}.Build()
	mockClient.On("DeleteCard", mock.Anything, mock.AnythingOfType("*proto.DeleteRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test deleting card
	err = client.DeleteCard(context.Background(), "card-1")
	assert.NoError(t, err)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

func TestClient_Sync(t *testing.T) {
	// Create a mock gRPC client
	mockClient := &mocks.KeeperServiceClient{}
	logger := zap.NewNop()

	client, err := NewClient(mockClient, logger)
	require.NoError(t, err)

	// Set up mock expectation
	now := time.Now()
	resp := proto.SyncResponse_builder{
		ServerTime: timestamppb.New(now),
	}.Build()

	mockClient.On("Sync", mock.Anything, mock.AnythingOfType("*proto.SyncRequest"), mock.Anything).Return(resp, nil)

	// Set a token to avoid file system dependency
	client.token = "test-token"

	// Test sync
	result, err := client.Sync(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.WithinDuration(t, now, result.GetServerTime().AsTime(), time.Second)

	// Verify mock was called
	mockClient.AssertExpectations(t)
}

// Helper functions for pointer creation
