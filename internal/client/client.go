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
	"google.golang.org/protobuf/types/known/emptypb"
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
	req := &proto.RegisterRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.grpcClient.Register(ctx, req)
	if err != nil {
		return "", fmt.Errorf("registration failed: %w", err)
	}

	return resp.Token, nil
}

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	req := &proto.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.grpcClient.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	return resp.Token, nil
}

func (c *Client) CreateCredential(ctx context.Context, title, login, password, meta string) (*proto.CredentialResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := &proto.CredentialRequest{
		Title:    title,
		Login:    login,
		Password: password,
		Meta:     meta,
	}

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

	req := &proto.UpdateCredentialRequest{
		Id:       id,
		Title:    title,
		Login:    login,
		Password: password,
		Meta:     meta,
		Version:  cred.Version,
	}

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

	req := &proto.DeleteRequest{
		Id: id,
	}

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

	req := &proto.GetRequest{
		Id: id,
	}

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

	req := &emptypb.Empty{}

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListCredentials(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list credentials from server, trying local storage", zap.Error(err))
		return c.localStorage.GetCredentials(), nil
	}

	return resp.Credentials, nil
}

func (c *Client) CreateTextData(ctx context.Context, title, data, meta string) (*proto.TextDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := &proto.TextDataRequest{
		Title: title,
		Data:  data,
		Meta:  meta,
	}

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

	req := &proto.UpdateTextDataRequest{
		Id:      id,
		Title:   title,
		Data:    data,
		Meta:    meta,
		Version: text.Version,
	}

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

	req := &proto.DeleteRequest{
		Id: id,
	}

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

	req := &proto.GetRequest{
		Id: id,
	}

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

	req := &emptypb.Empty{}

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListTextData(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list text data from server, trying local storage", zap.Error(err))
		return c.localStorage.GetTextData(), nil
	}

	return resp.TextData, nil
}

func (c *Client) CreateBinaryData(ctx context.Context, title string, data []byte, meta string) (*proto.BinaryDataResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	if len(data) > 10*1024*1024 {
		return nil, fmt.Errorf("data too large (max 10MB)")
	}

	req := &proto.BinaryDataRequest{
		Title: title,
		Data:  data,
		Meta:  meta,
	}

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

	req := &proto.UpdateBinaryDataRequest{
		Id:      id,
		Title:   title,
		Data:    data,
		Meta:    meta,
		Version: binary.Version,
	}

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

	req := &proto.DeleteRequest{
		Id: id,
	}

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

	req := &proto.GetRequest{
		Id: id,
	}

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

	req := &emptypb.Empty{}

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListBinaryData(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list binary data from server, trying local storage", zap.Error(err))
		return c.localStorage.GetBinaryData(), nil
	}

	return resp.BinaryData, nil
}

func (c *Client) CreateCard(ctx context.Context, title, cardNumber, cardHolder, expiry, cvv, meta string) (*proto.CardResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := &proto.CardRequest{
		Title:      title,
		CardNumber: cardNumber,
		CardHolder: cardHolder,
		Expiry:     expiry,
		Cvv:        cvv,
		Meta:       meta,
	}

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

	req := &proto.UpdateCardRequest{
		Id:         id,
		Title:      title,
		CardNumber: cardNumber,
		CardHolder: cardHolder,
		Expiry:     expiry,
		Cvv:        cvv,
		Meta:       meta,
		Version:    card.Version,
	}

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

	req := &proto.DeleteRequest{
		Id: id,
	}

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

	req := &proto.GetRequest{
		Id: id,
	}

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

	req := &emptypb.Empty{}

	ctx = c.withAuth(ctx)
	resp, err := c.grpcClient.ListCards(ctx, req)
	if err != nil {

		c.logger.Warn("failed to list cards from server, trying local storage", zap.Error(err))
		return c.localStorage.GetCards(), nil
	}

	return resp.Cards, nil
}

func (c *Client) Sync(ctx context.Context) (*proto.SyncResponse, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	req := &proto.SyncRequest{
		LastSync: timestamppb.New(time.Time{}),
	}

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
