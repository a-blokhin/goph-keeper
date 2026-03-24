package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	tokenFile = ".gophkeeper_token"
)

type Client struct {
	grpcClient   proto.KeeperServiceClient
	logger       *zap.Logger
	token        string
	localStorage *LocalStorage
}

func NewClient(grpcClient proto.KeeperServiceClient, logger *zap.Logger) (*Client, error) {
	localStorage, err := NewLocalStorage(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create local storage: %w", err)
	}

	return &Client{
		grpcClient:   grpcClient,
		logger:       logger,
		localStorage: localStorage,
	}, nil
}

func (c *Client) loadToken() error {
	if c.token != "" {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	tokenPath := filepath.Join(homeDir, tokenFile)
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to read token file: %w", err)
	}

	c.token = string(tokenBytes)
	return nil
}

func (c *Client) SaveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	tokenPath := filepath.Join(homeDir, tokenFile)
	if err := os.WriteFile(tokenPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	c.token = token
	return nil
}

func (c *Client) Register(ctx context.Context, email, password string) (string, error) {
	req := proto.RegisterRequest_builder{
		Email:    protobuf.String(email),
		Password: protobuf.String(password),
	}.Build()

	resp, err := c.grpcClient.Register(ctx, req)
	if err != nil {
		return "", fmt.Errorf("registration failed: %w", err)
	}

	return resp.GetToken(), nil
}

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	req := proto.LoginRequest_builder{
		Email:    protobuf.String(email),
		Password: protobuf.String(password),
	}.Build()

	resp, err := c.grpcClient.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	return resp.GetToken(), nil
}

func (c *Client) CreateCredential(ctx context.Context, title, login, password, meta string) (*proto.CredentialResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.CredentialRequest_builder{
		Title:    protobuf.String(title),
		Login:    protobuf.String(login),
		Password: protobuf.String(password),
		Meta:     protobuf.String(meta),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.CreateCredential(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	if err := c.localStorage.AddCredential(resp); err != nil {
		c.logger.Warn("failed to add credential to local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) UpdateCredential(ctx context.Context, id, title, login, password, meta string) (*proto.CredentialResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	cred, err := c.GetCredential(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	req := proto.UpdateCredentialRequest_builder{
		Id:       protobuf.String(id),
		Title:    protobuf.String(title),
		Login:    protobuf.String(login),
		Password: protobuf.String(password),
		Meta:     protobuf.String(meta),
		Version:  protobuf.Int32(cred.GetVersion()),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.UpdateCredential(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update credential: %w", err)
	}

	if err := c.localStorage.AddCredential(resp); err != nil {
		c.logger.Warn("failed to update credential in local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) DeleteCredential(ctx context.Context, id string) error {
	if err := c.loadToken(); err != nil {
		return err
	}

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	_, err := c.grpcClient.DeleteCredential(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	if err := c.localStorage.RemoveCredential(id); err != nil {
		c.logger.Warn("failed to remove credential from local storage", zap.Error(err))
	}

	return nil
}

func (c *Client) GetCredential(ctx context.Context, id string) (*proto.CredentialResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.GetRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.GetCredential(ctx, req)
	if err != nil {

		c.logger.Warn("failed to get credential from server, trying local storage", zap.Error(err))
		if localCred := c.localStorage.GetCredential(id); localCred != nil {
			return localCred, nil
		}
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	return resp, nil
}

func (c *Client) ListCredentials(ctx context.Context) ([]*proto.CredentialResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.ListCredentialsRequest_builder{}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListCredentials(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list credentials from server, trying local storage", zap.Error(err))
		return c.localStorage.GetCredentials(), nil
	}

	return resp.GetCredentials(), nil
}

func (c *Client) CreateTextData(ctx context.Context, title, data, meta string) (*proto.TextDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.TextDataRequest_builder{
		Title: protobuf.String(title),
		Data:  protobuf.String(data),
		Meta:  protobuf.String(meta),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.CreateTextData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create text data: %w", err)
	}

	if err := c.localStorage.AddTextData(resp); err != nil {
		c.logger.Warn("failed to add text data to local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) UpdateTextData(ctx context.Context, id, title, data, meta string) (*proto.TextDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	text, err := c.GetTextData(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get text data: %w", err)
	}

	req := proto.UpdateTextDataRequest_builder{
		Id:      protobuf.String(id),
		Title:   protobuf.String(title),
		Data:    protobuf.String(data),
		Meta:    protobuf.String(meta),
		Version: protobuf.Int32(text.GetVersion()),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.UpdateTextData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update text data: %w", err)
	}

	if err := c.localStorage.AddTextData(resp); err != nil {
		c.logger.Warn("failed to update text data in local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) DeleteTextData(ctx context.Context, id string) error {
	if err := c.loadToken(); err != nil {
		return err
	}

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	_, err := c.grpcClient.DeleteTextData(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete text data: %w", err)
	}

	if err := c.localStorage.RemoveTextData(id); err != nil {
		c.logger.Warn("failed to remove text data from local storage", zap.Error(err))
	}

	return nil
}

func (c *Client) GetTextData(ctx context.Context, id string) (*proto.TextDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.GetRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.GetTextData(ctx, req)
	if err != nil {

		c.logger.Warn("failed to get text data from server, trying local storage", zap.Error(err))
		if localText := c.localStorage.GetTextDataById(id); localText != nil {
			return localText, nil
		}
		return nil, fmt.Errorf("failed to get text data: %w", err)
	}

	return resp, nil
}

func (c *Client) ListTextData(ctx context.Context) ([]*proto.TextDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.ListTextDataRequest_builder{}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListTextData(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list text data from server, trying local storage", zap.Error(err))
		return c.localStorage.GetTextData(), nil
	}

	return resp.GetTextData(), nil
}

func (c *Client) CreateBinaryData(ctx context.Context, title string, data []byte, meta string) (*proto.BinaryDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	if len(data) > 10*1024*1024 {
		return nil, fmt.Errorf("data too large (max 10MB)")
	}

	req := proto.BinaryDataRequest_builder{
		Title: protobuf.String(title),
		Data:  data,
		Meta:  protobuf.String(meta),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.CreateBinaryData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create binary data: %w", err)
	}

	if err := c.localStorage.AddBinaryData(resp); err != nil {
		c.logger.Warn("failed to add binary data to local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) UpdateBinaryData(ctx context.Context, id, title string, data []byte, meta string) (*proto.BinaryDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	if len(data) > 10*1024*1024 {
		return nil, fmt.Errorf("data too large (max 10MB)")
	}

	binary, err := c.GetBinaryData(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get binary data: %w", err)
	}

	req := proto.UpdateBinaryDataRequest_builder{
		Id:      protobuf.String(id),
		Title:   protobuf.String(title),
		Data:    data,
		Meta:    protobuf.String(meta),
		Version: protobuf.Int32(binary.GetVersion()),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.UpdateBinaryData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update binary data: %w", err)
	}

	if err := c.localStorage.AddBinaryData(resp); err != nil {
		c.logger.Warn("failed to update binary data in local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) DeleteBinaryData(ctx context.Context, id string) error {
	if err := c.loadToken(); err != nil {
		return err
	}

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	_, err := c.grpcClient.DeleteBinaryData(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete binary data: %w", err)
	}

	if err := c.localStorage.RemoveBinaryData(id); err != nil {
		c.logger.Warn("failed to remove binary data from local storage", zap.Error(err))
	}

	return nil
}

func (c *Client) GetBinaryData(ctx context.Context, id string) (*proto.BinaryDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.GetRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.GetBinaryData(ctx, req)
	if err != nil {

		c.logger.Warn("failed to get binary data from server, trying local storage", zap.Error(err))
		if localBinary := c.localStorage.GetBinaryDataById(id); localBinary != nil {
			return localBinary, nil
		}
		return nil, fmt.Errorf("failed to get binary data: %w", err)
	}

	return resp, nil
}

func (c *Client) ListBinaryData(ctx context.Context) ([]*proto.BinaryDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.ListBinaryDataRequest_builder{}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListBinaryData(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list binary data from server, trying local storage", zap.Error(err))
		return c.localStorage.GetBinaryData(), nil
	}

	return resp.GetBinaryData(), nil
}

func (c *Client) CreateCard(ctx context.Context, title, cardNumber, cardHolder, expiry, cvv, meta string) (*proto.CardResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.CardRequest_builder{
		Title:      protobuf.String(title),
		CardNumber: protobuf.String(cardNumber),
		CardHolder: protobuf.String(cardHolder),
		Expiry:     protobuf.String(expiry),
		Cvv:        protobuf.String(cvv),
		Meta:       protobuf.String(meta),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.CreateCard(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	if err := c.localStorage.AddCard(resp); err != nil {
		c.logger.Warn("failed to add card to local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) UpdateCard(ctx context.Context, id, title, cardNumber, cardHolder, expiry, cvv, meta string) (*proto.CardResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	card, err := c.GetCard(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	req := proto.UpdateCardRequest_builder{
		Id:         protobuf.String(id),
		Title:      protobuf.String(title),
		CardNumber: protobuf.String(cardNumber),
		CardHolder: protobuf.String(cardHolder),
		Expiry:     protobuf.String(expiry),
		Cvv:        protobuf.String(cvv),
		Meta:       protobuf.String(meta),
		Version:    protobuf.Int32(card.GetVersion()),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.UpdateCard(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	if err := c.localStorage.AddCard(resp); err != nil {
		c.logger.Warn("failed to update card in local storage", zap.Error(err))
	}

	return resp, nil
}

func (c *Client) DeleteCard(ctx context.Context, id string) error {
	if err := c.loadToken(); err != nil {
		return err
	}

	req := proto.DeleteRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	_, err := c.grpcClient.DeleteCard(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	if err := c.localStorage.RemoveCard(id); err != nil {
		c.logger.Warn("failed to remove card from local storage", zap.Error(err))
	}

	return nil
}

func (c *Client) GetCard(ctx context.Context, id string) (*proto.CardResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.GetRequest_builder{
		Id: protobuf.String(id),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.GetCard(ctx, req)
	if err != nil {

		c.logger.Warn("failed to get card from server, trying local storage", zap.Error(err))
		if localCard := c.localStorage.GetCard(id); localCard != nil {
			return localCard, nil
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return resp, nil
}

func (c *Client) ListCards(ctx context.Context) ([]*proto.CardResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.ListCardsRequest_builder{}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListCards(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list cards from server, trying local storage", zap.Error(err))
		return c.localStorage.GetCards(), nil
	}

	return resp.GetCards(), nil
}

func (c *Client) Sync(ctx context.Context) (*proto.SyncResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := proto.SyncRequest_builder{
		LastSync: timestamppb.New(time.Time{}),
	}.Build()

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.Sync(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sync failed: %w", err)
	}

	return resp, nil
}

func (c *Client) withAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.token)
}
