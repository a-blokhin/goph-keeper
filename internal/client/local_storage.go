package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"go.uber.org/zap"
)

const (
	localDataFile = ".gophkeeper_data"
)

type LocalStorage struct {
	logger *zap.Logger
	data   *LocalData
	mu     sync.RWMutex
}

type SerializableLocalData struct {
	Credentials []map[string]interface{} `json:"credentials"`
	TextData    []map[string]interface{} `json:"text_data"`
	BinaryData  []map[string]interface{} `json:"binary_data"`
	Cards       []map[string]interface{} `json:"cards"`
	LastSync    time.Time                `json:"last_sync"`
}

type LocalData struct {
	Credentials []*proto.CredentialResponse
	TextData    []*proto.TextDataResponse
	BinaryData  []*proto.BinaryDataResponse
	Cards       []*proto.CardResponse
	LastSync    time.Time
}

func NewLocalStorage(logger *zap.Logger) (*LocalStorage, error) {
	storage := &LocalStorage{
		logger: logger,
		data: &LocalData{
			Credentials: []*proto.CredentialResponse{},
			TextData:    []*proto.TextDataResponse{},
			BinaryData:  []*proto.BinaryDataResponse{},
			Cards:       []*proto.CardResponse{},
			LastSync:    time.Time{},
		},
	}

	if err := storage.load(); err != nil {
		logger.Warn("Failed to load local data, starting fresh", zap.Error(err))
	}

	return storage, nil
}

func (s *LocalStorage) load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dataPath := filepath.Join(homeDir, localDataFile)
	dataBytes, err := os.ReadFile(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No local data yet
		}
		return fmt.Errorf("failed to read local data file: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var serializableData SerializableLocalData
	if err := json.Unmarshal(dataBytes, &serializableData); err != nil {
		return fmt.Errorf("failed to unmarshal local data: %w", err)
	}

	s.data = s.deserializeData(&serializableData)

	return nil
}

func (s *LocalStorage) serializeData() *SerializableLocalData {
	serializable := &SerializableLocalData{
		LastSync: s.data.LastSync,
	}

	// Serialize credentials
	for _, cred := range s.data.Credentials {
		credMap := map[string]interface{}{
			"id":       cred.GetId(),
			"title":    cred.GetTitle(),
			"login":    cred.GetLogin(),
			"password": cred.GetPassword(),
			"meta":     cred.GetMeta(),
			"version":  cred.GetVersion(),
		}
		serializable.Credentials = append(serializable.Credentials, credMap)
	}

	// Serialize text data
	for _, text := range s.data.TextData {
		textMap := map[string]interface{}{
			"id":      text.GetId(),
			"title":   text.GetTitle(),
			"data":    text.GetData(),
			"meta":    text.GetMeta(),
			"version": text.GetVersion(),
		}
		serializable.TextData = append(serializable.TextData, textMap)
	}

	// Serialize binary data
	for _, binary := range s.data.BinaryData {
		binaryMap := map[string]interface{}{
			"id":      binary.GetId(),
			"title":   binary.GetTitle(),
			"data":    binary.GetData(),
			"meta":    binary.GetMeta(),
			"version": binary.GetVersion(),
		}
		serializable.BinaryData = append(serializable.BinaryData, binaryMap)
	}

	// Serialize cards
	for _, card := range s.data.Cards {
		cardMap := map[string]interface{}{
			"id":          card.GetId(),
			"title":       card.GetTitle(),
			"card_number": card.GetCardNumber(),
			"card_holder": card.GetCardHolder(),
			"expiry":      card.GetExpiry(),
			"cvv":         card.GetCvv(),
			"meta":        card.GetMeta(),
			"version":     card.GetVersion(),
		}
		serializable.Cards = append(serializable.Cards, cardMap)
	}

	return serializable
}

func (s *LocalStorage) deserializeData(serializable *SerializableLocalData) *LocalData {
	data := &LocalData{
		LastSync: serializable.LastSync,
	}

	// Deserialize credentials
	for _, credMap := range serializable.Credentials {
		cred := proto.CredentialResponse_builder{
			Id:       stringPtr(credMap["id"].(string)),
			Title:    stringPtr(credMap["title"].(string)),
			Login:    stringPtr(credMap["login"].(string)),
			Password: stringPtr(credMap["password"].(string)),
			Meta:     stringPtr(credMap["meta"].(string)),
			Version:  int32Ptr(int32(credMap["version"].(float64))),
		}.Build()
		data.Credentials = append(data.Credentials, cred)
	}

	// Deserialize text data
	for _, textMap := range serializable.TextData {
		text := proto.TextDataResponse_builder{
			Id:      stringPtr(textMap["id"].(string)),
			Title:   stringPtr(textMap["title"].(string)),
			Data:    stringPtr(textMap["data"].(string)),
			Meta:    stringPtr(textMap["meta"].(string)),
			Version: int32Ptr(int32(textMap["version"].(float64))),
		}.Build()
		data.TextData = append(data.TextData, text)
	}

	// Deserialize binary data
	for _, binaryMap := range serializable.BinaryData {
		binary := proto.BinaryDataResponse_builder{
			Id:      stringPtr(binaryMap["id"].(string)),
			Title:   stringPtr(binaryMap["title"].(string)),
			Data:    []byte(binaryMap["data"].(string)),
			Meta:    stringPtr(binaryMap["meta"].(string)),
			Version: int32Ptr(int32(binaryMap["version"].(float64))),
		}.Build()
		data.BinaryData = append(data.BinaryData, binary)
	}

	// Deserialize cards
	for _, cardMap := range serializable.Cards {
		card := proto.CardResponse_builder{
			Id:         stringPtr(cardMap["id"].(string)),
			Title:      stringPtr(cardMap["title"].(string)),
			CardNumber: stringPtr(cardMap["card_number"].(string)),
			CardHolder: stringPtr(cardMap["card_holder"].(string)),
			Expiry:     stringPtr(cardMap["expiry"].(string)),
			Cvv:        stringPtr(cardMap["cvv"].(string)),
			Meta:       stringPtr(cardMap["meta"].(string)),
			Version:    int32Ptr(int32(cardMap["version"].(float64))),
		}.Build()
		data.Cards = append(data.Cards, card)
	}

	return data
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func (s *LocalStorage) save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dataPath := filepath.Join(homeDir, localDataFile)

	serializableData := s.serializeData()
	dataBytes, err := json.MarshalIndent(serializableData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal local data: %w", err)
	}

	if err := os.WriteFile(dataPath, dataBytes, 0600); err != nil {
		return fmt.Errorf("failed to write local data file: %w", err)
	}

	return nil
}

func (s *LocalStorage) UpdateFromSync(syncResp *proto.SyncResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.Credentials = syncResp.GetCredentials()
	s.data.TextData = syncResp.GetTextData()
	s.data.BinaryData = syncResp.GetBinaryData()
	s.data.Cards = syncResp.GetCards()
	s.data.LastSync = syncResp.GetServerTime().AsTime()

	return s.save()
}

func (s *LocalStorage) GetCredentials() []*proto.CredentialResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data.Credentials
}

func (s *LocalStorage) GetTextData() []*proto.TextDataResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data.TextData
}

func (s *LocalStorage) GetBinaryData() []*proto.BinaryDataResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data.BinaryData
}

func (s *LocalStorage) GetCards() []*proto.CardResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data.Cards
}

func (s *LocalStorage) GetLastSync() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data.LastSync
}

func (s *LocalStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = &LocalData{
		Credentials: []*proto.CredentialResponse{},
		TextData:    []*proto.TextDataResponse{},
		BinaryData:  []*proto.BinaryDataResponse{},
		Cards:       []*proto.CardResponse{},
		LastSync:    time.Time{},
	}

	return s.save()
}

func (s *LocalStorage) AddCredential(cred *proto.CredentialResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing credential with same ID if exists
	for i, existing := range s.data.Credentials {
		if existing.GetId() == cred.GetId() {
			s.data.Credentials = append(s.data.Credentials[:i], s.data.Credentials[i+1:]...)
			break
		}
	}

	s.data.Credentials = append(s.data.Credentials, cred)
	return s.save()
}

func (s *LocalStorage) AddTextData(text *proto.TextDataResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing text data with same ID if exists
	for i, existing := range s.data.TextData {
		if existing.GetId() == text.GetId() {
			s.data.TextData = append(s.data.TextData[:i], s.data.TextData[i+1:]...)
			break
		}
	}

	s.data.TextData = append(s.data.TextData, text)
	return s.save()
}

func (s *LocalStorage) AddBinaryData(binary *proto.BinaryDataResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing binary data with same ID if exists
	for i, existing := range s.data.BinaryData {
		if existing.GetId() == binary.GetId() {
			s.data.BinaryData = append(s.data.BinaryData[:i], s.data.BinaryData[i+1:]...)
			break
		}
	}

	s.data.BinaryData = append(s.data.BinaryData, binary)
	return s.save()
}

func (s *LocalStorage) AddCard(card *proto.CardResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing card with same ID if exists
	for i, existing := range s.data.Cards {
		if existing.GetId() == card.GetId() {
			s.data.Cards = append(s.data.Cards[:i], s.data.Cards[i+1:]...)
			break
		}
	}

	s.data.Cards = append(s.data.Cards, card)
	return s.save()
}

func (s *LocalStorage) RemoveCredential(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, cred := range s.data.Credentials {
		if cred.GetId() == id {
			s.data.Credentials = append(s.data.Credentials[:i], s.data.Credentials[i+1:]...)
			return s.save()
		}
	}

	return nil
}

func (s *LocalStorage) RemoveTextData(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, text := range s.data.TextData {
		if text.GetId() == id {
			s.data.TextData = append(s.data.TextData[:i], s.data.TextData[i+1:]...)
			return s.save()
		}
	}

	return nil
}

func (s *LocalStorage) RemoveBinaryData(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, binary := range s.data.BinaryData {
		if binary.GetId() == id {
			s.data.BinaryData = append(s.data.BinaryData[:i], s.data.BinaryData[i+1:]...)
			return s.save()
		}
	}

	return nil
}

func (s *LocalStorage) RemoveCard(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, card := range s.data.Cards {
		if card.GetId() == id {
			s.data.Cards = append(s.data.Cards[:i], s.data.Cards[i+1:]...)
			return s.save()
		}
	}

	return nil
}

func (s *LocalStorage) GetCredential(id string) *proto.CredentialResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, cred := range s.data.Credentials {
		if cred.GetId() == id {
			return cred
		}
	}

	return nil
}

func (s *LocalStorage) GetTextDataById(id string) *proto.TextDataResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, text := range s.data.TextData {
		if text.GetId() == id {
			return text
		}
	}

	return nil
}

func (s *LocalStorage) GetBinaryDataById(id string) *proto.BinaryDataResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, binary := range s.data.BinaryData {
		if binary.GetId() == id {
			return binary
		}
	}

	return nil
}

func (s *LocalStorage) GetCard(id string) *proto.CardResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, card := range s.data.Cards {
		if card.GetId() == id {
			return card
		}
	}

	return nil
}
