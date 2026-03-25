package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/repository/mocks"
	encryptionmocks "github.com/a-blokhin/goph-keeper/internal/service/encryption/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestSyncUsecase_Execute_Success(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"

	// Mock data
	expectedCredentials := []*model.Credential{
		{ID: "cred-1", UserID: userID, Title: "Test Credential", Login: "user1", PasswordEncrypted: "encrypted-pass1"},
	}
	expectedTextData := []*model.TextData{
		{ID: "text-1", UserID: userID, Title: "Test Text", DataEncrypted: "encrypted-text-data"},
	}
	expectedBinaryData := []*model.BinaryData{
		{ID: "binary-1", UserID: userID, Title: "Test Binary", DataEncrypted: []byte("encrypted-binary-data")},
	}
	expectedCards := []*model.Card{
		{ID: "card-1", UserID: userID, Title: "Test Card", CardNumberEncrypted: "encrypted-number", CardHolderEncrypted: "encrypted-holder", ExpiryEncrypted: "encrypted-expiry", CVVEncrypted: "encrypted-cvv"},
	}

	// Setup expectations
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCredentials, nil)
	mockTextDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedTextData, nil)
	mockBinaryDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedBinaryData, nil)
	mockCardRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCards, nil)

	mockEncryptionService.On("DecryptCredential", expectedCredentials[0]).Return(nil)
	mockEncryptionService.On("DecryptTextData", expectedTextData[0]).Return(nil)
	mockEncryptionService.On("DecryptBinaryData", expectedBinaryData[0]).Return(nil)
	mockEncryptionService.On("DecryptCardData", expectedCards[0]).Return(nil)

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCredentials, result.Credentials)
	assert.Equal(t, expectedTextData, result.TextData)
	assert.Equal(t, expectedBinaryData, result.BinaryData)
	assert.Equal(t, expectedCards, result.Cards)
	assert.WithinDuration(t, time.Now(), result.LastSync, time.Second)

	// Verify all expectations were met
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
	mockBinaryDataRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockEncryptionService.AssertExpectations(t)
}

func TestSyncUsecase_Execute_CredentialRepoError(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"
	expectedError := errors.New("database error")

	// Setup expectations
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return(nil, expectedError)

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	// Verify only credential repo was called
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertNotCalled(t, "GetByUserID")
	mockBinaryDataRepo.AssertNotCalled(t, "GetByUserID")
	mockCardRepo.AssertNotCalled(t, "GetByUserID")
	mockEncryptionService.AssertNotCalled(t, "DecryptCredential")
}

func TestSyncUsecase_Execute_TextDataRepoError(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"
	expectedError := errors.New("database error")
	expectedCredentials := []*model.Credential{
		{ID: "cred-1", UserID: userID, Title: "Test Credential", Login: "user1", PasswordEncrypted: "encrypted-pass1"},
	}

	// Setup expectations
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCredentials, nil)
	mockTextDataRepo.On("GetByUserID", mock.Anything, userID).Return(nil, expectedError)
	// No decryption calls should be made when text data repo fails

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
	mockBinaryDataRepo.AssertNotCalled(t, "GetByUserID")
	mockCardRepo.AssertNotCalled(t, "GetByUserID")
	mockEncryptionService.AssertExpectations(t)
}

func TestSyncUsecase_Execute_BinaryDataRepoError(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"
	expectedError := errors.New("database error")
	expectedCredentials := []*model.Credential{
		{ID: "cred-1", UserID: userID, Title: "Test Credential", Login: "user1", PasswordEncrypted: "encrypted-pass1"},
	}
	expectedTextData := []*model.TextData{
		{ID: "text-1", UserID: userID, Title: "Test Text", DataEncrypted: "encrypted-text-data"},
	}

	// Setup expectations
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCredentials, nil)
	mockTextDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedTextData, nil)
	mockBinaryDataRepo.On("GetByUserID", mock.Anything, userID).Return(nil, expectedError)

	// No decryption calls should be made when binary data repo fails

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
	mockBinaryDataRepo.AssertExpectations(t)
	mockCardRepo.AssertNotCalled(t, "GetByUserID")
	mockEncryptionService.AssertExpectations(t)
}

func TestSyncUsecase_Execute_CardRepoError(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"
	expectedError := errors.New("database error")
	expectedCredentials := []*model.Credential{
		{ID: "cred-1", UserID: userID, Title: "Test Credential", Login: "user1", PasswordEncrypted: "encrypted-pass1"},
	}
	expectedTextData := []*model.TextData{
		{ID: "text-1", UserID: userID, Title: "Test Text", DataEncrypted: "encrypted-text-data"},
	}
	expectedBinaryData := []*model.BinaryData{
		{ID: "binary-1", UserID: userID, Title: "Test Binary", DataEncrypted: []byte("encrypted-binary-data")},
	}

	// Setup expectations
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCredentials, nil)
	mockTextDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedTextData, nil)
	mockBinaryDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedBinaryData, nil)
	mockCardRepo.On("GetByUserID", mock.Anything, userID).Return(nil, expectedError)

	// No decryption calls should be made when card repo fails

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
	mockBinaryDataRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockEncryptionService.AssertExpectations(t)
}

func TestSyncUsecase_Execute_DecryptionError(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"
	expectedError := errors.New("decryption failed")
	expectedCredentials := []*model.Credential{
		{ID: "cred-1", UserID: userID, Title: "Test Credential", Login: "user1", PasswordEncrypted: "encrypted-pass1"},
	}
	expectedTextData := []*model.TextData{
		{ID: "text-1", UserID: userID, Title: "Test Text", DataEncrypted: "encrypted-text-data"},
	}
	expectedBinaryData := []*model.BinaryData{
		{ID: "binary-1", UserID: userID, Title: "Test Binary", DataEncrypted: []byte("encrypted-binary-data")},
	}
	expectedCards := []*model.Card{
		{ID: "card-1", UserID: userID, Title: "Test Card", CardNumberEncrypted: "encrypted-number", CardHolderEncrypted: "encrypted-holder", ExpiryEncrypted: "encrypted-expiry", CVVEncrypted: "encrypted-cvv"},
	}

	// Setup expectations
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCredentials, nil)
	mockTextDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedTextData, nil)
	mockBinaryDataRepo.On("GetByUserID", mock.Anything, userID).Return(expectedBinaryData, nil)
	mockCardRepo.On("GetByUserID", mock.Anything, userID).Return(expectedCards, nil)

	mockEncryptionService.On("DecryptCredential", expectedCredentials[0]).Return(expectedError)

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
	mockBinaryDataRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	mockEncryptionService.AssertExpectations(t)
}

func TestSyncUsecase_Execute_EmptyData(t *testing.T) {
	// Setup
	mockCredentialRepo := &mocks.CredentialRepository{}
	mockTextDataRepo := &mocks.TextDataRepository{}
	mockBinaryDataRepo := &mocks.BinaryDataRepository{}
	mockCardRepo := &mocks.CardRepository{}
	mockEncryptionService := &encryptionmocks.EncryptionService{}

	logger := zap.NewNop()

	usecase := New(
		mockCredentialRepo,
		mockTextDataRepo,
		mockBinaryDataRepo,
		mockCardRepo,
		mockEncryptionService,
		logger,
	)

	userID := "test-user-id"

	// Setup expectations - return empty slices
	mockCredentialRepo.On("GetByUserID", mock.Anything, userID).Return([]*model.Credential{}, nil)
	mockTextDataRepo.On("GetByUserID", mock.Anything, userID).Return([]*model.TextData{}, nil)
	mockBinaryDataRepo.On("GetByUserID", mock.Anything, userID).Return([]*model.BinaryData{}, nil)
	mockCardRepo.On("GetByUserID", mock.Anything, userID).Return([]*model.Card{}, nil)

	// Execute
	ctx := context.Background()
	result, err := usecase.Execute(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Credentials)
	assert.Empty(t, result.TextData)
	assert.Empty(t, result.BinaryData)
	assert.Empty(t, result.Cards)
	assert.WithinDuration(t, time.Now(), result.LastSync, time.Second)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
	mockBinaryDataRepo.AssertExpectations(t)
	mockCardRepo.AssertExpectations(t)
	// No decryption calls should be made for empty data
	mockEncryptionService.AssertNotCalled(t, "DecryptCredential")
	mockEncryptionService.AssertNotCalled(t, "DecryptTextData")
	mockEncryptionService.AssertNotCalled(t, "DecryptBinaryData")
	mockEncryptionService.AssertNotCalled(t, "DecryptCardData")
}
