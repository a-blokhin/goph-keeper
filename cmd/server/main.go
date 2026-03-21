package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/di/app"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	addr          = flag.String("addr", ":50051", "The server address")
	dsn           = flag.String("dsn", "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable", "Database connection string")
	jwtSecret     = flag.String("jwt-secret", "your-secret-key", "JWT secret key")
	encryptionKey = flag.String("encryption-key", "01234567890123456789012345678901", "32-byte encryption key for AES-GCM")
	tlsCert       = flag.String("tls-cert", "", "TLS certificate file path")
	tlsKey        = flag.String("tls-key", "", "TLS key file path")
	enableTLS     = flag.Bool("enable-tls", false, "Enable TLS")
)

func main() {
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		*addr = envAddr
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		*dsn = envDSN
	}
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		*jwtSecret = envSecret
	}
	if envEncryptionKey := os.Getenv("ENCRYPTION_KEY"); envEncryptionKey != "" {
		*encryptionKey = envEncryptionKey
	}
	if envCert := os.Getenv("TLS_CERT"); envCert != "" {
		*tlsCert = envCert
	}
	if envKey := os.Getenv("TLS_KEY"); envKey != "" {
		*tlsKey = envKey
	}
	if envTLS := os.Getenv("ENABLE_TLS"); envTLS != "" {
		*enableTLS = envTLS == "true" || envTLS == "1"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting GophKeeper server",
		zap.String("addr", *addr),
		zap.Bool("tls", *enableTLS),
	)

	ctx := context.Background()

	container, err := app.NewContainer(ctx, logger, *dsn, *jwtSecret, *encryptionKey)
	if err != nil {
		logger.Fatal("Failed to create container", zap.Error(err))
	}
	defer container.Close()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	var serverOpts []grpc.ServerOption

	if *enableTLS {
		if *tlsCert == "" || *tlsKey == "" {
			logger.Fatal("TLS certificate and key files are required when TLS is enabled")
		}

		creds, err := credentials.NewServerTLSFromFile(*tlsCert, *tlsKey)
		if err != nil {
			logger.Fatal("Failed to load TLS credentials", zap.Error(err))
		}

		serverOpts = append(serverOpts, grpc.Creds(creds))
		logger.Info("TLS enabled")
	}

	grpcServer := grpc.NewServer(serverOpts...)
	proto.RegisterKeeperServiceServer(grpcServer, container.Server)

	go func() {
		logger.Info("Server started", zap.String("addr", *addr))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
}
