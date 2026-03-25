package client

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewLocalStorage(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
	assert.NotNil(t, storage.data)
}

func TestLocalStorage_AddCredential(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add a credential
	cred := proto.CredentialResponse_builder{
		Id:       strPtr("cred-1"),
		Title:    strPtr("Test Credential"),
		Login:    strPtr("testuser"),
		Password: strPtr("testpass"),
		Meta:     strPtr("test meta"),
		Version:  i32Ptr(1),
	}.Build()

	err = storage.AddCredential(cred)
	assert.NoError(t, err)

	// Verify it was added
	creds := storage.GetCredentials()
	assert.Len(t, creds, 1)
	assert.Equal(t, "cred-1", creds[0].GetId())
	assert.Equal(t, "Test Credential", creds[0].GetTitle())

	// Add another credential with the same ID - should replace
	cred2 := proto.CredentialResponse_builder{
		Id:       strPtr("cred-1"),
		Title:    strPtr("Updated Credential"),
		Login:    strPtr("updateduser"),
		Password: strPtr("updatedpass"),
		Meta:     strPtr("updated meta"),
		Version:  i32Ptr(2),
	}.Build()

	err = storage.AddCredential(cred2)
	assert.NoError(t, err)

	// Verify it was replaced
	creds = storage.GetCredentials()
	assert.Len(t, creds, 1)
	assert.Equal(t, "cred-1", creds[0].GetId())
	assert.Equal(t, "Updated Credential", creds[0].GetTitle())
	assert.Equal(t, int32(2), creds[0].GetVersion())

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_AddTextData(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add text data
	text := proto.TextDataResponse_builder{
		Id:      strPtr("text-1"),
		Title:   strPtr("Test Text"),
		Data:    strPtr("test data content"),
		Meta:    strPtr("test meta"),
		Version: i32Ptr(1),
	}.Build()

	err = storage.AddTextData(text)
	assert.NoError(t, err)

	// Verify it was added
	textData := storage.GetTextData()
	assert.Len(t, textData, 1)
	assert.Equal(t, "text-1", textData[0].GetId())
	assert.Equal(t, "Test Text", textData[0].GetTitle())

	// Add another text data with the same ID - should replace
	text2 := proto.TextDataResponse_builder{
		Id:      strPtr("text-1"),
		Title:   strPtr("Updated Text"),
		Data:    strPtr("updated data content"),
		Meta:    strPtr("updated meta"),
		Version: i32Ptr(2),
	}.Build()

	err = storage.AddTextData(text2)
	assert.NoError(t, err)

	// Verify it was replaced
	textData = storage.GetTextData()
	assert.Len(t, textData, 1)
	assert.Equal(t, "text-1", textData[0].GetId())
	assert.Equal(t, "Updated Text", textData[0].GetTitle())
	assert.Equal(t, int32(2), textData[0].GetVersion())

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_AddBinaryData(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add binary data
	binary := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-1"),
		Title:   strPtr("Test Binary"),
		Data:    []byte("test binary content"),
		Meta:    strPtr("test meta"),
		Version: i32Ptr(1),
	}.Build()

	err = storage.AddBinaryData(binary)
	assert.NoError(t, err)

	// Verify it was added
	binaryData := storage.GetBinaryData()
	assert.Len(t, binaryData, 1)
	assert.Equal(t, "binary-1", binaryData[0].GetId())
	assert.Equal(t, "Test Binary", binaryData[0].GetTitle())

	// Add another binary data with the same ID - should replace
	binary2 := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-1"),
		Title:   strPtr("Updated Binary"),
		Data:    []byte("updated binary content"),
		Meta:    strPtr("updated meta"),
		Version: i32Ptr(2),
	}.Build()

	err = storage.AddBinaryData(binary2)
	assert.NoError(t, err)

	// Verify it was replaced
	binaryData = storage.GetBinaryData()
	assert.Len(t, binaryData, 1)
	assert.Equal(t, "binary-1", binaryData[0].GetId())
	assert.Equal(t, "Updated Binary", binaryData[0].GetTitle())
	assert.Equal(t, int32(2), binaryData[0].GetVersion())

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_AddCard(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add a card
	card := proto.CardResponse_builder{
		Id:         strPtr("card-1"),
		Title:      strPtr("Test Card"),
		CardNumber: strPtr("1234567890123456"),
		CardHolder: strPtr("Test User"),
		Expiry:     strPtr("12/25"),
		Cvv:        strPtr("123"),
		Meta:       strPtr("test meta"),
		Version:    i32Ptr(1),
	}.Build()

	err = storage.AddCard(card)
	assert.NoError(t, err)

	// Verify it was added
	cards := storage.GetCards()
	assert.Len(t, cards, 1)
	assert.Equal(t, "card-1", cards[0].GetId())
	assert.Equal(t, "Test Card", cards[0].GetTitle())

	// Add another card with the same ID - should replace
	card2 := proto.CardResponse_builder{
		Id:         strPtr("card-1"),
		Title:      strPtr("Updated Card"),
		CardNumber: strPtr("6543210987654321"),
		CardHolder: strPtr("Updated User"),
		Expiry:     strPtr("11/26"),
		Cvv:        strPtr("456"),
		Meta:       strPtr("updated meta"),
		Version:    i32Ptr(2),
	}.Build()

	err = storage.AddCard(card2)
	assert.NoError(t, err)

	// Verify it was replaced
	cards = storage.GetCards()
	assert.Len(t, cards, 1)
	assert.Equal(t, "card-1", cards[0].GetId())
	assert.Equal(t, "Updated Card", cards[0].GetTitle())
	assert.Equal(t, int32(2), cards[0].GetVersion())

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_RemoveCredential(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add credentials
	cred1 := proto.CredentialResponse_builder{
		Id:       strPtr("cred-1"),
		Title:    strPtr("Test Credential 1"),
		Login:    strPtr("testuser1"),
		Password: strPtr("testpass1"),
		Meta:     strPtr("test meta 1"),
		Version:  i32Ptr(1),
	}.Build()

	cred2 := proto.CredentialResponse_builder{
		Id:       strPtr("cred-2"),
		Title:    strPtr("Test Credential 2"),
		Login:    strPtr("testuser2"),
		Password: strPtr("testpass2"),
		Meta:     strPtr("test meta 2"),
		Version:  i32Ptr(1),
	}.Build()

	err = storage.AddCredential(cred1)
	require.NoError(t, err)

	err = storage.AddCredential(cred2)
	require.NoError(t, err)

	// Verify both are there
	creds := storage.GetCredentials()
	assert.Len(t, creds, 2)

	// Remove one
	err = storage.RemoveCredential("cred-1")
	assert.NoError(t, err)

	// Verify only one remains
	creds = storage.GetCredentials()
	assert.Len(t, creds, 1)
	assert.Equal(t, "cred-2", creds[0].GetId())

	// Try to remove non-existent credential - should not error
	err = storage.RemoveCredential("non-existent")
	assert.NoError(t, err)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_RemoveTextData(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add text data
	text1 := proto.TextDataResponse_builder{
		Id:      strPtr("text-1"),
		Title:   strPtr("Test Text 1"),
		Data:    strPtr("test data content 1"),
		Meta:    strPtr("test meta 1"),
		Version: i32Ptr(1),
	}.Build()

	text2 := proto.TextDataResponse_builder{
		Id:      strPtr("text-2"),
		Title:   strPtr("Test Text 2"),
		Data:    strPtr("test data content 2"),
		Meta:    strPtr("test meta 2"),
		Version: i32Ptr(1),
	}.Build()

	err = storage.AddTextData(text1)
	require.NoError(t, err)

	err = storage.AddTextData(text2)
	require.NoError(t, err)

	// Verify both are there
	textData := storage.GetTextData()
	assert.Len(t, textData, 2)

	// Remove one
	err = storage.RemoveTextData("text-1")
	assert.NoError(t, err)

	// Verify only one remains
	textData = storage.GetTextData()
	assert.Len(t, textData, 1)
	assert.Equal(t, "text-2", textData[0].GetId())

	// Try to remove non-existent text data - should not error
	err = storage.RemoveTextData("non-existent")
	assert.NoError(t, err)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_RemoveBinaryData(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add binary data
	binary1 := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-1"),
		Title:   strPtr("Test Binary 1"),
		Data:    []byte("test binary content 1"),
		Meta:    strPtr("test meta 1"),
		Version: i32Ptr(1),
	}.Build()

	binary2 := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-2"),
		Title:   strPtr("Test Binary 2"),
		Data:    []byte("test binary content 2"),
		Meta:    strPtr("test meta 2"),
		Version: i32Ptr(1),
	}.Build()

	err = storage.AddBinaryData(binary1)
	require.NoError(t, err)

	err = storage.AddBinaryData(binary2)
	require.NoError(t, err)

	// Verify both are there
	binaryData := storage.GetBinaryData()
	assert.Len(t, binaryData, 2)

	// Remove one
	err = storage.RemoveBinaryData("binary-1")
	assert.NoError(t, err)

	// Verify only one remains
	binaryData = storage.GetBinaryData()
	assert.Len(t, binaryData, 1)
	assert.Equal(t, "binary-2", binaryData[0].GetId())

	// Try to remove non-existent binary data - should not error
	err = storage.RemoveBinaryData("non-existent")
	assert.NoError(t, err)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_RemoveCard(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add cards
	card1 := proto.CardResponse_builder{
		Id:         strPtr("card-1"),
		Title:      strPtr("Test Card 1"),
		CardNumber: strPtr("1234567890123456"),
		CardHolder: strPtr("Test User 1"),
		Expiry:     strPtr("12/25"),
		Cvv:        strPtr("123"),
		Meta:       strPtr("test meta 1"),
		Version:    i32Ptr(1),
	}.Build()

	card2 := proto.CardResponse_builder{
		Id:         strPtr("card-2"),
		Title:      strPtr("Test Card 2"),
		CardNumber: strPtr("6543210987654321"),
		CardHolder: strPtr("Test User 2"),
		Expiry:     strPtr("11/26"),
		Cvv:        strPtr("456"),
		Meta:       strPtr("test meta 2"),
		Version:    i32Ptr(1),
	}.Build()

	err = storage.AddCard(card1)
	require.NoError(t, err)

	err = storage.AddCard(card2)
	require.NoError(t, err)

	// Verify both are there
	cards := storage.GetCards()
	assert.Len(t, cards, 2)

	// Remove one
	err = storage.RemoveCard("card-1")
	assert.NoError(t, err)

	// Verify only one remains
	cards = storage.GetCards()
	assert.Len(t, cards, 1)
	assert.Equal(t, "card-2", cards[0].GetId())

	// Try to remove non-existent card - should not error
	err = storage.RemoveCard("non-existent")
	assert.NoError(t, err)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_GetCredential(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add credentials
	cred1 := proto.CredentialResponse_builder{
		Id:       strPtr("cred-1"),
		Title:    strPtr("Test Credential 1"),
		Login:    strPtr("testuser1"),
		Password: strPtr("testpass1"),
		Meta:     strPtr("test meta 1"),
		Version:  i32Ptr(1),
	}.Build()

	cred2 := proto.CredentialResponse_builder{
		Id:       strPtr("cred-2"),
		Title:    strPtr("Test Credential 2"),
		Login:    strPtr("testuser2"),
		Password: strPtr("testpass2"),
		Meta:     strPtr("test meta 2"),
		Version:  i32Ptr(1),
	}.Build()

	err = storage.AddCredential(cred1)
	require.NoError(t, err)

	err = storage.AddCredential(cred2)
	require.NoError(t, err)

	// Get existing credential
	result := storage.GetCredential("cred-1")
	assert.NotNil(t, result)
	assert.Equal(t, "cred-1", result.GetId())
	assert.Equal(t, "Test Credential 1", result.GetTitle())

	// Get non-existent credential
	result = storage.GetCredential("non-existent")
	assert.Nil(t, result)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_GetTextDataById(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add text data
	text1 := proto.TextDataResponse_builder{
		Id:      strPtr("text-1"),
		Title:   strPtr("Test Text 1"),
		Data:    strPtr("test data content 1"),
		Meta:    strPtr("test meta 1"),
		Version: i32Ptr(1),
	}.Build()

	text2 := proto.TextDataResponse_builder{
		Id:      strPtr("text-2"),
		Title:   strPtr("Test Text 2"),
		Data:    strPtr("test data content 2"),
		Meta:    strPtr("test meta 2"),
		Version: i32Ptr(1),
	}.Build()

	err = storage.AddTextData(text1)
	require.NoError(t, err)

	err = storage.AddTextData(text2)
	require.NoError(t, err)

	// Get existing text data
	result := storage.GetTextDataById("text-1")
	assert.NotNil(t, result)
	assert.Equal(t, "text-1", result.GetId())
	assert.Equal(t, "Test Text 1", result.GetTitle())

	// Get non-existent text data
	result = storage.GetTextDataById("non-existent")
	assert.Nil(t, result)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_GetBinaryDataById(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add binary data
	binary1 := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-1"),
		Title:   strPtr("Test Binary 1"),
		Data:    []byte("test binary content 1"),
		Meta:    strPtr("test meta 1"),
		Version: i32Ptr(1),
	}.Build()

	binary2 := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-2"),
		Title:   strPtr("Test Binary 2"),
		Data:    []byte("test binary content 2"),
		Meta:    strPtr("test meta 2"),
		Version: i32Ptr(1),
	}.Build()

	err = storage.AddBinaryData(binary1)
	require.NoError(t, err)

	err = storage.AddBinaryData(binary2)
	require.NoError(t, err)

	// Get existing binary data
	result := storage.GetBinaryDataById("binary-1")
	assert.NotNil(t, result)
	assert.Equal(t, "binary-1", result.GetId())
	assert.Equal(t, "Test Binary 1", result.GetTitle())

	// Get non-existent binary data
	result = storage.GetBinaryDataById("non-existent")
	assert.Nil(t, result)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_GetCard(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add cards
	card1 := proto.CardResponse_builder{
		Id:         strPtr("card-1"),
		Title:      strPtr("Test Card 1"),
		CardNumber: strPtr("1234567890123456"),
		CardHolder: strPtr("Test User 1"),
		Expiry:     strPtr("12/25"),
		Cvv:        strPtr("123"),
		Meta:       strPtr("test meta 1"),
		Version:    i32Ptr(1),
	}.Build()

	card2 := proto.CardResponse_builder{
		Id:         strPtr("card-2"),
		Title:      strPtr("Test Card 2"),
		CardNumber: strPtr("6543210987654321"),
		CardHolder: strPtr("Test User 2"),
		Expiry:     strPtr("11/26"),
		Cvv:        strPtr("456"),
		Meta:       strPtr("test meta 2"),
		Version:    i32Ptr(1),
	}.Build()

	err = storage.AddCard(card1)
	require.NoError(t, err)

	err = storage.AddCard(card2)
	require.NoError(t, err)

	// Get existing card
	result := storage.GetCard("card-1")
	assert.NotNil(t, result)
	assert.Equal(t, "card-1", result.GetId())
	assert.Equal(t, "Test Card 1", result.GetTitle())

	// Get non-existent card
	result = storage.GetCard("non-existent")
	assert.Nil(t, result)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_Clear(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add some data
	cred := proto.CredentialResponse_builder{
		Id:       strPtr("cred-1"),
		Title:    strPtr("Test Credential"),
		Login:    strPtr("testuser"),
		Password: strPtr("testpass"),
		Meta:     strPtr("test meta"),
		Version:  i32Ptr(1),
	}.Build()

	text := proto.TextDataResponse_builder{
		Id:      strPtr("text-1"),
		Title:   strPtr("Test Text"),
		Data:    strPtr("test data content"),
		Meta:    strPtr("test meta"),
		Version: i32Ptr(1),
	}.Build()

	binary := proto.BinaryDataResponse_builder{
		Id:      strPtr("binary-1"),
		Title:   strPtr("Test Binary"),
		Data:    []byte("test binary content"),
		Meta:    strPtr("test meta"),
		Version: i32Ptr(1),
	}.Build()

	card := proto.CardResponse_builder{
		Id:         strPtr("card-1"),
		Title:      strPtr("Test Card"),
		CardNumber: strPtr("1234567890123456"),
		CardHolder: strPtr("Test User"),
		Expiry:     strPtr("12/25"),
		Cvv:        strPtr("123"),
		Meta:       strPtr("test meta"),
		Version:    i32Ptr(1),
	}.Build()

	err = storage.AddCredential(cred)
	require.NoError(t, err)

	err = storage.AddTextData(text)
	require.NoError(t, err)

	err = storage.AddBinaryData(binary)
	require.NoError(t, err)

	err = storage.AddCard(card)
	require.NoError(t, err)

	// Verify data was added
	assert.Len(t, storage.GetCredentials(), 1)
	assert.Len(t, storage.GetTextData(), 1)
	assert.Len(t, storage.GetBinaryData(), 1)
	assert.Len(t, storage.GetCards(), 1)

	// Clear all data
	err = storage.Clear()
	assert.NoError(t, err)

	// Verify all data was cleared
	assert.Len(t, storage.GetCredentials(), 0)
	assert.Len(t, storage.GetTextData(), 0)
	assert.Len(t, storage.GetBinaryData(), 0)
	assert.Len(t, storage.GetCards(), 0)
	assert.Equal(t, time.Time{}, storage.GetLastSync())

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_UpdateFromSync(t *testing.T) {
	logger := zap.NewNop()
	storage, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Create sync response with data
	now := time.Now()
	syncResp := proto.SyncResponse_builder{
		ServerTime: timestamppb.New(now),
		Credentials: []*proto.CredentialResponse{
			proto.CredentialResponse_builder{
				Id:       strPtr("cred-1"),
				Title:    strPtr("Synced Credential"),
				Login:    strPtr("synceduser"),
				Password: strPtr("syncedpass"),
				Meta:     strPtr("synced meta"),
				Version:  i32Ptr(1),
			}.Build(),
		},
		TextData: []*proto.TextDataResponse{
			proto.TextDataResponse_builder{
				Id:      strPtr("text-1"),
				Title:   strPtr("Synced Text"),
				Data:    strPtr("synced data content"),
				Meta:    strPtr("synced meta"),
				Version: i32Ptr(1),
			}.Build(),
		},
		BinaryData: []*proto.BinaryDataResponse{
			proto.BinaryDataResponse_builder{
				Id:      strPtr("binary-1"),
				Title:   strPtr("Synced Binary"),
				Data:    []byte("synced binary content"),
				Meta:    strPtr("synced meta"),
				Version: i32Ptr(1),
			}.Build(),
		},
		Cards: []*proto.CardResponse{
			proto.CardResponse_builder{
				Id:         strPtr("card-1"),
				Title:      strPtr("Synced Card"),
				CardNumber: strPtr("1234567890123456"),
				CardHolder: strPtr("Synced User"),
				Expiry:     strPtr("12/25"),
				Cvv:        strPtr("123"),
				Meta:       strPtr("synced meta"),
				Version:    i32Ptr(1),
			}.Build(),
		},
	}.Build()

	// Update from sync
	err = storage.UpdateFromSync(syncResp)
	assert.NoError(t, err)

	// Verify data was updated
	assert.Len(t, storage.GetCredentials(), 1)
	assert.Len(t, storage.GetTextData(), 1)
	assert.Len(t, storage.GetBinaryData(), 1)
	assert.Len(t, storage.GetCards(), 1)
	assert.WithinDuration(t, now, storage.GetLastSync(), time.Second)

	// Clean up
	_ = storage.Clear()
}

func TestLocalStorage_LoadSave(t *testing.T) {
	logger := zap.NewNop()

	// Create a temporary storage and add some data
	storage1, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Add some data
	cred := proto.CredentialResponse_builder{
		Id:       strPtr("cred-1"),
		Title:    strPtr("Test Credential"),
		Login:    strPtr("testuser"),
		Password: strPtr("testpass"),
		Meta:     strPtr("test meta"),
		Version:  i32Ptr(1),
	}.Build()

	err = storage1.AddCredential(cred)
	require.NoError(t, err)

	// Get the data file path
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	dataPath := filepath.Join(homeDir, localDataFile)

	// Create a new storage instance which should load the data
	storage2, err := NewLocalStorage(logger)
	require.NoError(t, err)

	// Verify data was loaded
	creds := storage2.GetCredentials()
	assert.Len(t, creds, 1)
	assert.Equal(t, "cred-1", creds[0].GetId())
	assert.Equal(t, "Test Credential", creds[0].GetTitle())

	// Clean up
	_ = storage1.Clear()
	_ = os.Remove(dataPath)
}

// Helper functions for pointer creation (local storage specific)
func strPtr(s string) *string {
	return &s
}

func i32Ptr(i int32) *int32 {
	return &i
}
