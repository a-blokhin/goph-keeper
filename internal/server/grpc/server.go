package grpc

import (
	"context"
	"errors"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/model"
	"github.com/a-blokhin/goph-keeper/internal/usecase/auth/login"
	"github.com/a-blokhin/goph-keeper/internal/usecase/auth/register"
	binaryDataCreate "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/create"
	binaryDataDelete "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/delete"
	binaryDataGet "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/get"
	binaryDataList "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/list"
	binaryDataUpdate "github.com/a-blokhin/goph-keeper/internal/usecase/binarydata/update"
	cardCreate "github.com/a-blokhin/goph-keeper/internal/usecase/card/create"
	cardDelete "github.com/a-blokhin/goph-keeper/internal/usecase/card/delete"
	cardGet "github.com/a-blokhin/goph-keeper/internal/usecase/card/get"
	cardList "github.com/a-blokhin/goph-keeper/internal/usecase/card/list"
	cardUpdate "github.com/a-blokhin/goph-keeper/internal/usecase/card/update"
	credentialCreate "github.com/a-blokhin/goph-keeper/internal/usecase/credential/create"
	credentialDelete "github.com/a-blokhin/goph-keeper/internal/usecase/credential/delete"
	credentialGet "github.com/a-blokhin/goph-keeper/internal/usecase/credential/get"
	credentialList "github.com/a-blokhin/goph-keeper/internal/usecase/credential/list"
	credentialUpdate "github.com/a-blokhin/goph-keeper/internal/usecase/credential/update"
	"github.com/a-blokhin/goph-keeper/internal/usecase/sync"
	textDataCreate "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/create"
	textDataDelete "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/delete"
	textDataGet "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/get"
	textDataList "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/list"
	textDataUpdate "github.com/a-blokhin/goph-keeper/internal/usecase/textdata/update"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedKeeperServiceServer
	logger          *zap.Logger
	jwtManager      *jwt.JWTManager
	registerUsecase register.RegisterUsecase
	loginUsecase    login.LoginUsecase

	credentialCreateUsecase credentialCreate.CreateCredentialUsecase
	credentialGetUsecase    credentialGet.GetCredentialUsecase
	credentialUpdateUsecase credentialUpdate.UpdateCredentialUsecase
	credentialDeleteUsecase credentialDelete.DeleteCredentialUsecase
	credentialListUsecase   credentialList.ListCredentialsUsecase

	textDataCreateUsecase textDataCreate.CreateTextDataUsecase
	textDataGetUsecase    textDataGet.GetTextDataUsecase
	textDataUpdateUsecase textDataUpdate.UpdateTextDataUsecase
	textDataDeleteUsecase textDataDelete.DeleteTextDataUsecase
	textDataListUsecase   textDataList.ListTextDataUsecase

	binaryDataCreateUsecase binaryDataCreate.CreateBinaryDataUsecase
	binaryDataGetUsecase    binaryDataGet.GetBinaryDataUsecase
	binaryDataUpdateUsecase binaryDataUpdate.UpdateBinaryDataUsecase
	binaryDataDeleteUsecase binaryDataDelete.DeleteBinaryDataUsecase
	binaryDataListUsecase   binaryDataList.ListBinaryDataUsecase

	cardCreateUsecase cardCreate.CreateCardUsecase
	cardGetUsecase    cardGet.GetCardUsecase
	cardUpdateUsecase cardUpdate.UpdateCardUsecase
	cardDeleteUsecase cardDelete.DeleteCardUsecase
	cardListUsecase   cardList.ListCardsUsecase

	syncUsecase sync.SyncUsecase
}

func NewServer(
	logger *zap.Logger,
	jwtManager *jwt.JWTManager,
	registerUsecase register.RegisterUsecase,
	loginUsecase login.LoginUsecase,
	credentialCreateUsecase credentialCreate.CreateCredentialUsecase,
	credentialGetUsecase credentialGet.GetCredentialUsecase,
	credentialUpdateUsecase credentialUpdate.UpdateCredentialUsecase,
	credentialDeleteUsecase credentialDelete.DeleteCredentialUsecase,
	credentialListUsecase credentialList.ListCredentialsUsecase,
	textDataCreateUsecase textDataCreate.CreateTextDataUsecase,
	textDataGetUsecase textDataGet.GetTextDataUsecase,
	textDataUpdateUsecase textDataUpdate.UpdateTextDataUsecase,
	textDataDeleteUsecase textDataDelete.DeleteTextDataUsecase,
	textDataListUsecase textDataList.ListTextDataUsecase,
	binaryDataCreateUsecase binaryDataCreate.CreateBinaryDataUsecase,
	binaryDataGetUsecase binaryDataGet.GetBinaryDataUsecase,
	binaryDataUpdateUsecase binaryDataUpdate.UpdateBinaryDataUsecase,
	binaryDataDeleteUsecase binaryDataDelete.DeleteBinaryDataUsecase,
	binaryDataListUsecase binaryDataList.ListBinaryDataUsecase,
	cardCreateUsecase cardCreate.CreateCardUsecase,
	cardGetUsecase cardGet.GetCardUsecase,
	cardUpdateUsecase cardUpdate.UpdateCardUsecase,
	cardDeleteUsecase cardDelete.DeleteCardUsecase,
	cardListUsecase cardList.ListCardsUsecase,
	syncUsecase sync.SyncUsecase,
) *Server {
	return &Server{
		logger:                  logger,
		jwtManager:              jwtManager,
		registerUsecase:         registerUsecase,
		loginUsecase:            loginUsecase,
		credentialCreateUsecase: credentialCreateUsecase,
		credentialGetUsecase:    credentialGetUsecase,
		credentialUpdateUsecase: credentialUpdateUsecase,
		credentialDeleteUsecase: credentialDeleteUsecase,
		credentialListUsecase:   credentialListUsecase,
		textDataCreateUsecase:   textDataCreateUsecase,
		textDataGetUsecase:      textDataGetUsecase,
		textDataUpdateUsecase:   textDataUpdateUsecase,
		textDataDeleteUsecase:   textDataDeleteUsecase,
		textDataListUsecase:     textDataListUsecase,
		binaryDataCreateUsecase: binaryDataCreateUsecase,
		binaryDataGetUsecase:    binaryDataGetUsecase,
		binaryDataUpdateUsecase: binaryDataUpdateUsecase,
		binaryDataDeleteUsecase: binaryDataDeleteUsecase,
		binaryDataListUsecase:   binaryDataListUsecase,
		cardCreateUsecase:       cardCreateUsecase,
		cardGetUsecase:          cardGetUsecase,
		cardUpdateUsecase:       cardUpdateUsecase,
		cardDeleteUsecase:       cardDeleteUsecase,
		cardListUsecase:         cardListUsecase,
		syncUsecase:             syncUsecase,
	}
}

func (s *Server) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	s.logger.Info("Register request", zap.String("email", req.GetEmail()))

	user, token, err := s.registerUsecase.Execute(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		s.logger.Error("Register failed", zap.Error(err))
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "registration failed")
	}

	s.logger.Info("User registered", zap.String("user_id", user.ID))

	return proto.RegisterResponse_builder{
		Token: protobuf.String(token),
	}.Build(), nil
}

func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	s.logger.Info("Login request", zap.String("email", req.GetEmail()))

	user, token, err := s.loginUsecase.Execute(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		s.logger.Error("Login failed", zap.Error(err))
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	s.logger.Info("User logged in", zap.String("user_id", user.ID))

	return proto.LoginResponse_builder{
		Token: protobuf.String(token),
	}.Build(), nil
}

func (s *Server) CreateCredential(ctx context.Context, req *proto.CredentialRequest) (*proto.CredentialResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("CreateCredential request", zap.String("user_id", userID), zap.String("title", req.GetTitle()))

	credential := &model.Credential{
		Title:             req.GetTitle(),
		Login:             req.GetLogin(),
		PasswordEncrypted: req.GetPassword(),
		Meta:              req.GetMeta(),
	}

	err = s.credentialCreateUsecase.Execute(ctx, userID, credential)
	if err != nil {
		s.logger.Error("CreateCredential failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create credential")
	}

	return s.credentialToProto(credential)
}

func (s *Server) GetCredential(ctx context.Context, req *proto.GetRequest) (*proto.CredentialResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("GetCredential request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	credential, err := s.credentialGetUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("GetCredential failed", zap.Error(err))
		if errors.Is(err, model.ErrCredentialNotFound) {
			return nil, status.Error(codes.NotFound, "credential not found")
		}
		return nil, status.Error(codes.Internal, "failed to get credential")
	}

	return s.credentialToProto(credential)
}

func (s *Server) UpdateCredential(ctx context.Context, req *proto.UpdateCredentialRequest) (*proto.CredentialResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("UpdateCredential request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	credential := &model.Credential{
		ID:                req.GetId(),
		Title:             req.GetTitle(),
		Login:             req.GetLogin(),
		PasswordEncrypted: req.GetPassword(),
		Meta:              req.GetMeta(),
		Version:           req.GetVersion(),
	}

	err = s.credentialUpdateUsecase.Execute(ctx, userID, credential)
	if err != nil {
		s.logger.Error("UpdateCredential failed", zap.Error(err))
		if errors.Is(err, model.ErrVersionConflict) {
			return nil, status.Error(codes.Aborted, "version conflict")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		if errors.Is(err, model.ErrCredentialNotFound) {
			return nil, status.Error(codes.NotFound, "credential not found")
		}
		return nil, status.Error(codes.Internal, "failed to update credential")
	}

	return s.credentialToProto(credential)
}

func (s *Server) DeleteCredential(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("DeleteCredential request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	err = s.credentialDeleteUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("DeleteCredential failed", zap.Error(err))
		if errors.Is(err, model.ErrCredentialNotFound) {
			return nil, status.Error(codes.NotFound, "credential not found")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		return nil, status.Error(codes.Internal, "failed to delete credential")
	}

	return proto.DeleteResponse_builder{}.Build(), nil
}

func (s *Server) ListCredentials(ctx context.Context, _ *proto.ListCredentialsRequest) (*proto.CredentialsListResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("ListCredentials request", zap.String("user_id", userID))

	credentials, err := s.credentialListUsecase.Execute(ctx, userID)
	if err != nil {
		s.logger.Error("ListCredentials failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list credentials")
	}

	protoCredentials := make([]*proto.CredentialResponse, len(credentials))
	for i, cred := range credentials {
		protoCred, err := s.credentialToProto(cred)
		if err != nil {
			s.logger.Error("Failed to convert credential to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert credential")
		}
		protoCredentials[i] = protoCred
	}

	return proto.CredentialsListResponse_builder{
		Credentials: protoCredentials,
	}.Build(), nil
}

func (s *Server) CreateTextData(ctx context.Context, req *proto.TextDataRequest) (*proto.TextDataResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("CreateTextData request", zap.String("user_id", userID), zap.String("title", req.GetTitle()))

	textData := &model.TextData{
		Title:         req.GetTitle(),
		DataEncrypted: req.GetData(),
		Meta:          req.GetMeta(),
	}

	err = s.textDataCreateUsecase.Execute(ctx, userID, textData)
	if err != nil {
		s.logger.Error("CreateTextData failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create text data")
	}

	return s.textDataToProto(textData)
}

func (s *Server) GetTextData(ctx context.Context, req *proto.GetRequest) (*proto.TextDataResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("GetTextData request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	textData, err := s.textDataGetUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("GetTextData failed", zap.Error(err))
		if errors.Is(err, model.ErrTextDataNotFound) {
			return nil, status.Error(codes.NotFound, "text data not found")
		}
		return nil, status.Error(codes.Internal, "failed to get text data")
	}

	return s.textDataToProto(textData)
}

func (s *Server) UpdateTextData(ctx context.Context, req *proto.UpdateTextDataRequest) (*proto.TextDataResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("UpdateTextData request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	textData := &model.TextData{
		ID:            req.GetId(),
		Title:         req.GetTitle(),
		DataEncrypted: req.GetData(),
		Meta:          req.GetMeta(),
		Version:       req.GetVersion(),
	}

	err = s.textDataUpdateUsecase.Execute(ctx, userID, textData)
	if err != nil {
		s.logger.Error("UpdateTextData failed", zap.Error(err))
		if errors.Is(err, model.ErrVersionConflict) {
			return nil, status.Error(codes.Aborted, "version conflict")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		if errors.Is(err, model.ErrTextDataNotFound) {
			return nil, status.Error(codes.NotFound, "text data not found")
		}
		return nil, status.Error(codes.Internal, "failed to update text data")
	}

	return s.textDataToProto(textData)
}

func (s *Server) DeleteTextData(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("DeleteTextData request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	err = s.textDataDeleteUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("DeleteTextData failed", zap.Error(err))
		if errors.Is(err, model.ErrTextDataNotFound) {
			return nil, status.Error(codes.NotFound, "text data not found")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		return nil, status.Error(codes.Internal, "failed to delete text data")
	}

	return proto.DeleteResponse_builder{}.Build(), nil
}

func (s *Server) ListTextData(ctx context.Context, _ *proto.ListTextDataRequest) (*proto.TextDataListResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("ListTextData request", zap.String("user_id", userID))

	textDataList, err := s.textDataListUsecase.Execute(ctx, userID)
	if err != nil {
		s.logger.Error("ListTextData failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list text data")
	}

	protoTextData := make([]*proto.TextDataResponse, len(textDataList))
	for i, td := range textDataList {
		protoTD, err := s.textDataToProto(td)
		if err != nil {
			s.logger.Error("Failed to convert text data to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert text data")
		}
		protoTextData[i] = protoTD
	}

	return proto.TextDataListResponse_builder{
		TextData: protoTextData,
	}.Build(), nil
}

func (s *Server) CreateBinaryData(ctx context.Context, req *proto.BinaryDataRequest) (*proto.BinaryDataResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("CreateBinaryData request", zap.String("user_id", userID), zap.String("title", req.GetTitle()))

	binaryData := &model.BinaryData{
		Title:         req.GetTitle(),
		DataEncrypted: req.GetData(),
		Meta:          req.GetMeta(),
	}

	err = s.binaryDataCreateUsecase.Execute(ctx, userID, binaryData)
	if err != nil {
		s.logger.Error("CreateBinaryData failed", zap.Error(err))
		if errors.Is(err, model.ErrBinaryDataTooLarge) {
			return nil, status.Error(codes.InvalidArgument, "data too large (max 10MB)")
		}
		return nil, status.Error(codes.Internal, "failed to create binary data")
	}

	return s.binaryDataToProto(binaryData)
}

func (s *Server) GetBinaryData(ctx context.Context, req *proto.GetRequest) (*proto.BinaryDataResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("GetBinaryData request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	binaryData, err := s.binaryDataGetUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("GetBinaryData failed", zap.Error(err))
		if errors.Is(err, model.ErrBinaryDataNotFound) {
			return nil, status.Error(codes.NotFound, "binary data not found")
		}
		return nil, status.Error(codes.Internal, "failed to get binary data")
	}

	return s.binaryDataToProto(binaryData)
}

func (s *Server) UpdateBinaryData(ctx context.Context, req *proto.UpdateBinaryDataRequest) (*proto.BinaryDataResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("UpdateBinaryData request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	binaryData := &model.BinaryData{
		ID:            req.GetId(),
		Title:         req.GetTitle(),
		DataEncrypted: req.GetData(),
		Meta:          req.GetMeta(),
		Version:       req.GetVersion(),
	}

	err = s.binaryDataUpdateUsecase.Execute(ctx, userID, binaryData)
	if err != nil {
		s.logger.Error("UpdateBinaryData failed", zap.Error(err))
		if errors.Is(err, model.ErrVersionConflict) {
			return nil, status.Error(codes.Aborted, "version conflict")
		}
		if errors.Is(err, model.ErrBinaryDataTooLarge) {
			return nil, status.Error(codes.InvalidArgument, "data too large (max 10MB)")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		if errors.Is(err, model.ErrBinaryDataNotFound) {
			return nil, status.Error(codes.NotFound, "binary data not found")
		}
		return nil, status.Error(codes.Internal, "failed to update binary data")
	}

	return s.binaryDataToProto(binaryData)
}

func (s *Server) DeleteBinaryData(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("DeleteBinaryData request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	err = s.binaryDataDeleteUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("DeleteBinaryData failed", zap.Error(err))
		if errors.Is(err, model.ErrBinaryDataNotFound) {
			return nil, status.Error(codes.NotFound, "binary data not found")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		return nil, status.Error(codes.Internal, "failed to delete binary data")
	}

	return proto.DeleteResponse_builder{}.Build(), nil
}

func (s *Server) ListBinaryData(ctx context.Context, _ *proto.ListBinaryDataRequest) (*proto.BinaryDataListResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("ListBinaryData request", zap.String("user_id", userID))

	binaryDataList, err := s.binaryDataListUsecase.Execute(ctx, userID)
	if err != nil {
		s.logger.Error("ListBinaryData failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list binary data")
	}

	protoBinaryData := make([]*proto.BinaryDataResponse, len(binaryDataList))
	for i, bd := range binaryDataList {
		protoBD, err := s.binaryDataToProto(bd)
		if err != nil {
			s.logger.Error("Failed to convert binary data to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert binary data")
		}
		protoBinaryData[i] = protoBD
	}

	return proto.BinaryDataListResponse_builder{
		BinaryData: protoBinaryData,
	}.Build(), nil
}

func (s *Server) CreateCard(ctx context.Context, req *proto.CardRequest) (*proto.CardResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("CreateCard request", zap.String("user_id", userID), zap.String("title", req.GetTitle()))

	card := &model.Card{
		Title:               req.GetTitle(),
		CardNumberEncrypted: req.GetCardNumber(),
		CardHolderEncrypted: req.GetCardHolder(),
		ExpiryEncrypted:     req.GetExpiry(),
		CVVEncrypted:        req.GetCvv(),
		Meta:                req.GetMeta(),
	}

	err = s.cardCreateUsecase.Execute(ctx, userID, card)
	if err != nil {
		s.logger.Error("CreateCard failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create card")
	}

	return s.cardToProto(card)
}

func (s *Server) GetCard(ctx context.Context, req *proto.GetRequest) (*proto.CardResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("GetCard request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	card, err := s.cardGetUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("GetCard failed", zap.Error(err))
		if errors.Is(err, model.ErrCardNotFound) {
			return nil, status.Error(codes.NotFound, "card not found")
		}
		return nil, status.Error(codes.Internal, "failed to get card")
	}

	return s.cardToProto(card)
}

func (s *Server) UpdateCard(ctx context.Context, req *proto.UpdateCardRequest) (*proto.CardResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("UpdateCard request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	card := &model.Card{
		ID:                  req.GetId(),
		Title:               req.GetTitle(),
		CardNumberEncrypted: req.GetCardNumber(),
		CardHolderEncrypted: req.GetCardHolder(),
		ExpiryEncrypted:     req.GetExpiry(),
		CVVEncrypted:        req.GetCvv(),
		Meta:                req.GetMeta(),
		Version:             req.GetVersion(),
	}

	err = s.cardUpdateUsecase.Execute(ctx, userID, card)
	if err != nil {
		s.logger.Error("UpdateCard failed", zap.Error(err))
		if errors.Is(err, model.ErrVersionConflict) {
			return nil, status.Error(codes.Aborted, "version conflict")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		if errors.Is(err, model.ErrCardNotFound) {
			return nil, status.Error(codes.NotFound, "card not found")
		}
		return nil, status.Error(codes.Internal, "failed to update card")
	}

	return s.cardToProto(card)
}

func (s *Server) DeleteCard(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("DeleteCard request", zap.String("user_id", userID), zap.String("id", req.GetId()))

	err = s.cardDeleteUsecase.Execute(ctx, userID, req.GetId())
	if err != nil {
		s.logger.Error("DeleteCard failed", zap.Error(err))
		if errors.Is(err, model.ErrCardNotFound) {
			return nil, status.Error(codes.NotFound, "card not found")
		}
		if errors.Is(err, model.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}
		return nil, status.Error(codes.Internal, "failed to delete card")
	}

	return proto.DeleteResponse_builder{}.Build(), nil
}

func (s *Server) ListCards(ctx context.Context, _ *proto.ListCardsRequest) (*proto.CardsListResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("ListCards request", zap.String("user_id", userID))

	cards, err := s.cardListUsecase.Execute(ctx, userID)
	if err != nil {
		s.logger.Error("ListCards failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list cards")
	}

	protoCards := make([]*proto.CardResponse, len(cards))
	for i, card := range cards {
		protoCard, err := s.cardToProto(card)
		if err != nil {
			s.logger.Error("Failed to convert card to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert card")
		}
		protoCards[i] = protoCard
	}

	return proto.CardsListResponse_builder{
		Cards: protoCards,
	}.Build(), nil
}

func (s *Server) Sync(ctx context.Context, req *proto.SyncRequest) (*proto.SyncResponse, error) {
	userID, err := s.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Sync request", zap.String("user_id", userID))

	resp, err := s.syncUsecase.Execute(ctx, userID)
	if err != nil {
		s.logger.Error("Sync failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "sync failed")
	}

	credentials := make([]*proto.CredentialResponse, len(resp.Credentials))
	for i, cred := range resp.Credentials {
		protoCred, err := s.credentialToProto(cred)
		if err != nil {
			s.logger.Error("Failed to convert credential to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert credential")
		}
		credentials[i] = protoCred
	}

	textData := make([]*proto.TextDataResponse, len(resp.TextData))
	for i, td := range resp.TextData {
		protoTD, err := s.textDataToProto(td)
		if err != nil {
			s.logger.Error("Failed to convert text data to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert text data")
		}
		textData[i] = protoTD
	}

	binaryData := make([]*proto.BinaryDataResponse, len(resp.BinaryData))
	for i, bd := range resp.BinaryData {
		protoBD, err := s.binaryDataToProto(bd)
		if err != nil {
			s.logger.Error("Failed to convert binary data to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert binary data")
		}
		binaryData[i] = protoBD
	}

	cards := make([]*proto.CardResponse, len(resp.Cards))
	for i, card := range resp.Cards {
		protoCard, err := s.cardToProto(card)
		if err != nil {
			s.logger.Error("Failed to convert card to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to convert card")
		}
		cards[i] = protoCard
	}

	return proto.SyncResponse_builder{
		Credentials: credentials,
		TextData:    textData,
		BinaryData:  binaryData,
		Cards:       cards,
		ServerTime:  timestamppb.New(resp.LastSync),
	}.Build(), nil
}

func (s *Server) extractUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	tokens := md["authorization"]
	if len(tokens) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization token")
	}

	token := tokens[0]
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims.UserID, nil
}

func (s *Server) credentialToProto(cred *model.Credential) (*proto.CredentialResponse, error) {
	version := int32(cred.Version)
	return proto.CredentialResponse_builder{
		Id:        protobuf.String(cred.ID),
		Title:     protobuf.String(cred.Title),
		Login:     protobuf.String(cred.Login),
		Password:  protobuf.String(cred.PasswordEncrypted),
		Meta:      protobuf.String(cred.Meta),
		Version:   protobuf.Int32(version),
		CreatedAt: timestamppb.New(cred.CreatedAt),
		UpdatedAt: timestamppb.New(cred.UpdatedAt),
	}.Build(), nil
}

func (s *Server) textDataToProto(td *model.TextData) (*proto.TextDataResponse, error) {
	version := int32(td.Version)
	return proto.TextDataResponse_builder{
		Id:        protobuf.String(td.ID),
		Title:     protobuf.String(td.Title),
		Data:      protobuf.String(td.DataEncrypted),
		Meta:      protobuf.String(td.Meta),
		Version:   protobuf.Int32(version),
		CreatedAt: timestamppb.New(td.CreatedAt),
		UpdatedAt: timestamppb.New(td.UpdatedAt),
	}.Build(), nil
}

func (s *Server) binaryDataToProto(bd *model.BinaryData) (*proto.BinaryDataResponse, error) {
	version := int32(bd.Version)
	return proto.BinaryDataResponse_builder{
		Id:        protobuf.String(bd.ID),
		Title:     protobuf.String(bd.Title),
		Data:      bd.DataEncrypted,
		Meta:      protobuf.String(bd.Meta),
		Version:   protobuf.Int32(version),
		CreatedAt: timestamppb.New(bd.CreatedAt),
		UpdatedAt: timestamppb.New(bd.UpdatedAt),
	}.Build(), nil
}

func (s *Server) cardToProto(card *model.Card) (*proto.CardResponse, error) {
	version := int32(card.Version)
	return proto.CardResponse_builder{
		Id:         protobuf.String(card.ID),
		Title:      protobuf.String(card.Title),
		CardNumber: protobuf.String(card.CardNumberEncrypted),
		CardHolder: protobuf.String(card.CardHolderEncrypted),
		Expiry:     protobuf.String(card.ExpiryEncrypted),
		Cvv:        protobuf.String(card.CVVEncrypted),
		Meta:       protobuf.String(card.Meta),
		Version:    protobuf.Int32(version),
		CreatedAt:  timestamppb.New(card.CreatedAt),
		UpdatedAt:  timestamppb.New(card.UpdatedAt),
	}.Build(), nil
}
