package app

import (
	"context"
	"fmt"

	"github.com/a-blokhin/goph-keeper/internal/crypto"
	"github.com/a-blokhin/goph-keeper/internal/jwt"
	"github.com/a-blokhin/goph-keeper/internal/migration"
	"github.com/a-blokhin/goph-keeper/internal/repository"
	"github.com/a-blokhin/goph-keeper/internal/repository/postgres"
	"github.com/a-blokhin/goph-keeper/internal/server/grpc"
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
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Container struct {
	Logger     *zap.Logger
	DBPool     *pgxpool.Pool
	JWTManager *jwt.JWTManager
	Encryptor  *crypto.Encryptor
	Server     *grpc.Server
}

func NewContainer(ctx context.Context, logger *zap.Logger, dsn string, jwtSecret string, encryptionKey string) (*Container, error) {
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	migrator := migration.New(logger, "./migrations")
	if err := migrator.Up(dsn); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	jwtManager := jwt.NewJWTManager(jwtSecret)

	encryptor, err := crypto.NewEncryptor([]byte(encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %w", err)
	}

	repos := initRepositories(dbPool)
	usecases := initUsecases(repos, jwtManager, logger)
	server := initServer(logger, jwtManager, encryptor, usecases)

	return &Container{
		Logger:     logger,
		DBPool:     dbPool,
		JWTManager: jwtManager,
		Encryptor:  encryptor,
		Server:     server,
	}, nil
}

func (c *Container) Close() {
	if c.DBPool != nil {
		c.DBPool.Close()
	}
}

type repositories struct {
	userRepo       repository.UserRepository
	credentialRepo repository.CredentialRepository
	textDataRepo   repository.TextDataRepository
	binaryDataRepo repository.BinaryDataRepository
	cardRepo       repository.CardRepository
}

func initRepositories(dbPool *pgxpool.Pool) *repositories {
	return &repositories{
		userRepo:       postgres.NewUserRepository(dbPool),
		credentialRepo: postgres.NewCredentialRepository(dbPool),
		textDataRepo:   postgres.NewTextDataRepository(dbPool),
		binaryDataRepo: postgres.NewBinaryDataRepository(dbPool),
		cardRepo:       postgres.NewCardRepository(dbPool),
	}
}

type usecases struct {
	registerUsecase         register.RegisterUsecase
	loginUsecase            login.LoginUsecase
	credentialCreateUsecase credentialCreate.CreateCredentialUsecase
	credentialGetUsecase    credentialGet.GetCredentialUsecase
	credentialUpdateUsecase credentialUpdate.UpdateCredentialUsecase
	credentialDeleteUsecase credentialDelete.DeleteCredentialUsecase
	credentialListUsecase   credentialList.ListCredentialsUsecase
	textDataCreateUsecase   textDataCreate.CreateTextDataUsecase
	textDataGetUsecase      textDataGet.GetTextDataUsecase
	textDataUpdateUsecase   textDataUpdate.UpdateTextDataUsecase
	textDataDeleteUsecase   textDataDelete.DeleteTextDataUsecase
	textDataListUsecase     textDataList.ListTextDataUsecase
	binaryDataCreateUsecase binaryDataCreate.CreateBinaryDataUsecase
	binaryDataGetUsecase    binaryDataGet.GetBinaryDataUsecase
	binaryDataUpdateUsecase binaryDataUpdate.UpdateBinaryDataUsecase
	binaryDataDeleteUsecase binaryDataDelete.DeleteBinaryDataUsecase
	binaryDataListUsecase   binaryDataList.ListBinaryDataUsecase
	cardCreateUsecase       cardCreate.CreateCardUsecase
	cardGetUsecase          cardGet.GetCardUsecase
	cardUpdateUsecase       cardUpdate.UpdateCardUsecase
	cardDeleteUsecase       cardDelete.DeleteCardUsecase
	cardListUsecase         cardList.ListCardsUsecase
	syncUsecase             sync.SyncUsecase
}

func initUsecases(repos *repositories, jwtManager *jwt.JWTManager, logger *zap.Logger) *usecases {
	registerUsecase := register.New(repos.userRepo, jwtManager, logger)
	loginUsecase := login.New(repos.userRepo, jwtManager, logger)

	credentialCreateUsecase := credentialCreate.New(repos.credentialRepo, logger)
	credentialGetUsecase := credentialGet.New(repos.credentialRepo)
	credentialUpdateUsecase := credentialUpdate.New(repos.credentialRepo, logger)
	credentialDeleteUsecase := credentialDelete.New(repos.credentialRepo, logger)
	credentialListUsecase := credentialList.New(repos.credentialRepo, logger)

	textDataCreateUsecase := textDataCreate.New(repos.textDataRepo, logger)
	textDataGetUsecase := textDataGet.New(repos.textDataRepo)
	textDataUpdateUsecase := textDataUpdate.New(repos.textDataRepo, logger)
	textDataDeleteUsecase := textDataDelete.New(repos.textDataRepo, logger)
	textDataListUsecase := textDataList.New(repos.textDataRepo, logger)

	binaryDataCreateUsecase := binaryDataCreate.New(repos.binaryDataRepo, logger)
	binaryDataGetUsecase := binaryDataGet.New(repos.binaryDataRepo)
	binaryDataUpdateUsecase := binaryDataUpdate.New(repos.binaryDataRepo, logger)
	binaryDataDeleteUsecase := binaryDataDelete.New(repos.binaryDataRepo, logger)
	binaryDataListUsecase := binaryDataList.New(repos.binaryDataRepo, logger)

	cardCreateUsecase := cardCreate.New(repos.cardRepo, logger)
	cardGetUsecase := cardGet.New(repos.cardRepo)
	cardUpdateUsecase := cardUpdate.New(repos.cardRepo, logger)
	cardDeleteUsecase := cardDelete.New(repos.cardRepo, logger)
	cardListUsecase := cardList.New(repos.cardRepo, logger)

	syncUsecase := sync.New(
		repos.credentialRepo,
		repos.textDataRepo,
		repos.binaryDataRepo,
		repos.cardRepo,
		logger,
	)

	return &usecases{
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

func initServer(logger *zap.Logger, jwtManager *jwt.JWTManager, encryptor *crypto.Encryptor, usecases *usecases) *grpc.Server {
	return grpc.NewServer(
		logger,
		jwtManager,
		encryptor,
		usecases.registerUsecase,
		usecases.loginUsecase,
		usecases.credentialCreateUsecase,
		usecases.credentialGetUsecase,
		usecases.credentialUpdateUsecase,
		usecases.credentialDeleteUsecase,
		usecases.credentialListUsecase,
		usecases.textDataCreateUsecase,
		usecases.textDataGetUsecase,
		usecases.textDataUpdateUsecase,
		usecases.textDataDeleteUsecase,
		usecases.textDataListUsecase,
		usecases.binaryDataCreateUsecase,
		usecases.binaryDataGetUsecase,
		usecases.binaryDataUpdateUsecase,
		usecases.binaryDataDeleteUsecase,
		usecases.binaryDataListUsecase,
		usecases.cardCreateUsecase,
		usecases.cardGetUsecase,
		usecases.cardUpdateUsecase,
		usecases.cardDeleteUsecase,
		usecases.cardListUsecase,
		usecases.syncUsecase,
	)
}
