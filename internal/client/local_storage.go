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

type LocalData struct {
	Credentials []*proto.CredentialResponse `json:"credentials"`
	TextData    []*proto.TextDataResponse   `json:"text_data"`
	BinaryData  []*proto.BinaryDataResponse `json:"binary_data"`
	Cards       []*proto.CardResponse       `json:"cards"`
	LastSync    time.Time                   `json:"last_sync"`
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

	if err := json.Unmarshal(dataBytes, &s.data); err != nil {
		return fmt.Errorf("failed to unmarshal local data: %w", err)
	}

	return nil
}

func (s *LocalStorage) save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dataPath := filepath.Join(homeDir, localDataFile)

	dataBytes, err := json.MarshalIndent(s.data, "", "  ")
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
