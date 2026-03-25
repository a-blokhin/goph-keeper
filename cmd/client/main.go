package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/a-blokhin/goph-keeper/api/proto"
	"github.com/a-blokhin/goph-keeper/internal/client"
	"go.uber.org/zap"
	"golang.org/x/term"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	version = "dev"

	serverAddr = flag.String("server", "localhost:50051", "Server address")
	enableTLS  = flag.Bool("tls", false, "Enable TLS")
	tlsCert    = flag.String("tls-cert", "", "TLS certificate file path (for testing)")
	username   = flag.String("username", "", "Username/email for register/login")
)

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	if len(flag.Args()) < 1 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Args()[0]
	args := flag.Args()[1:]

	ctx := context.Background()

	conn, err := createConnection(*serverAddr, *enableTLS, *tlsCert)
	if err != nil {
		logger.Fatal("Failed to connect to server", zap.Error(err))
	}
	defer conn.Close()

	grpcClient := proto.NewKeeperServiceClient(conn)
	cli, err := client.NewClient(grpcClient, logger)
	if err != nil {
		logger.Fatal("Failed to create client", zap.Error(err))
	}

	switch command {
	case "register":
		handleRegister(ctx, cli, args)
	case "login":
		handleLogin(ctx, cli, args)
	case "credential":
		handleCredential(ctx, cli, args)
	case "text":
		handleText(ctx, cli, args)
	case "binary":
		handleBinary(ctx, cli, args)
	case "card":
		handleCard(ctx, cli, args)
	case "sync":
		handleSync(ctx, cli, args)
	case "version":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func createConnection(addr string, enableTLS bool, tlsCert string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if enableTLS {
		if tlsCert != "" {
			certPool := x509.NewCertPool()
			certBytes, err := os.ReadFile(tlsCert)
			if err != nil {
				return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
			}
			if !certPool.AppendCertsFromPEM(certBytes) {
				return nil, fmt.Errorf("failed to parse TLS certificate")
			}
			creds := credentials.NewTLS(&tls.Config{
				RootCAs: certPool,
			})
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			creds := credentials.NewTLS(&tls.Config{})
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return grpc.NewClient(addr, opts...)
}

func handleRegister(ctx context.Context, cli *client.Client, args []string) {
	email := *username
	var password string

	if email == "" {
		if len(args) >= 1 {
			email = args[0]
		} else {
			fmt.Print("Enter email: ")
			fmt.Scanln(&email)
		}
	}

	if email == "" {
		fmt.Fprintln(os.Stderr, "Error: email is required")
		fmt.Fprintln(os.Stderr, "Usage: register <email>")
		fmt.Fprintln(os.Stderr, "       register --username <email>")
		os.Exit(1)
	}

	password = readPassword("Enter password: ")
	if password == "" {
		fmt.Fprintln(os.Stderr, "Error: password is required")
		os.Exit(1)
	}

	token, err := cli.Register(ctx, email, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Registration failed: %v\n", err)
		os.Exit(1)
	}

	if err := cli.SaveToken(token); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save token: %v\n", err)
	}

	fmt.Println("Registration successful!")
}

func handleLogin(ctx context.Context, cli *client.Client, args []string) {
	email := *username
	var password string

	if email == "" {
		if len(args) >= 1 {
			email = args[0]
		} else {
			fmt.Print("Enter email: ")
			fmt.Scanln(&email)
		}
	}

	if email == "" {
		fmt.Fprintln(os.Stderr, "Error: email is required")
		fmt.Fprintln(os.Stderr, "Usage: login <email>")
		fmt.Fprintln(os.Stderr, "       login --username <email>")
		os.Exit(1)
	}

	password = readPassword("Enter password: ")
	if password == "" {
		fmt.Fprintln(os.Stderr, "Error: password is required")
		os.Exit(1)
	}

	token, err := cli.Login(ctx, email, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
		os.Exit(1)
	}

	if err := cli.SaveToken(token); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save token: %v\n", err)
	}

	fmt.Println("Login successful!")
}

func readPassword(prompt string) string {
	fmt.Print(prompt)

	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		log.Printf("Error reading password: %v", err)
		return ""
	}

	return strings.TrimSpace(string(passwordBytes))
}

func handleCredential(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: credential <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  create <title> <login> <password> [meta]")
		fmt.Println("  update <id> <title> <login> <password> [meta]")
		fmt.Println("  delete <id>")
		fmt.Println("  list")
		os.Exit(1)
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "create":
		handleCreateCredential(ctx, cli, subargs)
	case "update":
		handleUpdateCredential(ctx, cli, subargs)
	case "delete":
		handleDeleteCredential(ctx, cli, subargs)
	case "list":
		handleListCredentials(ctx, cli)
	default:
		fmt.Printf("Unknown credential command: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleCreateCredential(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: credential create <title> <login> <password> [meta]")
		os.Exit(1)
	}

	title := args[0]
	login := args[1]
	password := args[2]
	meta := ""
	if len(args) > 3 {
		meta = args[3]
	}

	cred, err := cli.CreateCredential(ctx, title, login, password, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create credential: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Credential created with ID: %s\n", cred.GetId())
}

func handleUpdateCredential(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: credential update <id> <title> <login> <password> [meta]")
		os.Exit(1)
	}

	id := args[0]
	title := args[1]
	login := args[2]
	password := args[3]
	meta := ""
	if len(args) > 4 {
		meta = args[4]
	}

	cred, err := cli.UpdateCredential(ctx, id, title, login, password, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update credential: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Credential updated: %s\n", cred.GetId())
}

func handleDeleteCredential(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: credential delete <id>")
		os.Exit(1)
	}

	id := args[0]

	if err := cli.DeleteCredential(ctx, id); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete credential: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Credential deleted")
}

func handleListCredentials(ctx context.Context, cli *client.Client) {
	creds, err := cli.ListCredentials(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list credentials: %v\n", err)
		os.Exit(1)
	}

	if len(creds) == 0 {
		fmt.Println("No credentials found")
		return
	}

	fmt.Println("Credentials:")
	for _, cred := range creds {
		fmt.Printf("  ID: %s, Title: %s, Login: %s\n", cred.GetId(), cred.GetTitle(), cred.GetLogin())
	}
}

func handleText(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: text <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  create <title> <data> [meta]")
		fmt.Println("  update <id> <title> <data> [meta]")
		fmt.Println("  delete <id>")
		fmt.Println("  list")
		os.Exit(1)
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "create":
		handleCreateText(ctx, cli, subargs)
	case "update":
		handleUpdateText(ctx, cli, subargs)
	case "delete":
		handleDeleteText(ctx, cli, subargs)
	case "list":
		handleListText(ctx, cli)
	default:
		fmt.Printf("Unknown text command: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleCreateText(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: text create <title> <data> [meta]")
		os.Exit(1)
	}

	title := args[0]
	data := args[1]
	meta := ""
	if len(args) > 2 {
		meta = args[2]
	}

	text, err := cli.CreateTextData(ctx, title, data, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create text data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Text data created with ID: %s\n", text.GetId())
}

func handleUpdateText(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: text update <id> <title> <data> [meta]")
		os.Exit(1)
	}

	id := args[0]
	title := args[1]
	data := args[2]
	meta := ""
	if len(args) > 3 {
		meta = args[3]
	}

	text, err := cli.UpdateTextData(ctx, id, title, data, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update text data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Text data updated: %s\n", text.GetId())
}

func handleDeleteText(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: text delete <id>")
		os.Exit(1)
	}

	id := args[0]

	if err := cli.DeleteTextData(ctx, id); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete text data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Text data deleted")
}

func handleListText(ctx context.Context, cli *client.Client) {
	texts, err := cli.ListTextData(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list text data: %v\n", err)
		os.Exit(1)
	}

	if len(texts) == 0 {
		fmt.Println("No text data found")
		return
	}

	fmt.Println("Text Data:")
	for _, text := range texts {
		fmt.Printf("  ID: %s, Title: %s\n", text.GetId(), text.GetTitle())
	}
}

func handleBinary(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: binary <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  create <title> <file> [meta]")
		fmt.Println("  update <id> <title> <file> [meta]")
		fmt.Println("  delete <id>")
		fmt.Println("  list")
		os.Exit(1)
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "create":
		handleCreateBinary(ctx, cli, subargs)
	case "update":
		handleUpdateBinary(ctx, cli, subargs)
	case "delete":
		handleDeleteBinary(ctx, cli, subargs)
	case "list":
		handleListBinary(ctx, cli)
	default:
		fmt.Printf("Unknown binary command: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleCreateBinary(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: binary create <title> <file> [meta]")
		os.Exit(1)
	}

	title := args[0]
	filePath := args[1]
	meta := ""
	if len(args) > 2 {
		meta = args[2]
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	binary, err := cli.CreateBinaryData(ctx, title, data, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create binary data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Binary data created with ID: %s\n", binary.GetId())
}

func handleUpdateBinary(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: binary update <id> <title> <file> [meta]")
		os.Exit(1)
	}

	id := args[0]
	title := args[1]
	filePath := args[2]
	meta := ""
	if len(args) > 3 {
		meta = args[3]
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	binary, err := cli.UpdateBinaryData(ctx, id, title, data, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update binary data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Binary data updated: %s\n", binary.GetId())
}

func handleDeleteBinary(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: binary delete <id>")
		os.Exit(1)
	}

	id := args[0]

	if err := cli.DeleteBinaryData(ctx, id); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete binary data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Binary data deleted")
}

func handleListBinary(ctx context.Context, cli *client.Client) {
	binaries, err := cli.ListBinaryData(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list binary data: %v\n", err)
		os.Exit(1)
	}

	if len(binaries) == 0 {
		fmt.Println("No binary data found")
		return
	}

	fmt.Println("Binary Data:")
	for _, binary := range binaries {
		fmt.Printf("  ID: %s, Title: %s, Size: %d bytes\n", binary.GetId(), binary.GetTitle(), len(binary.GetData()))
	}
}

func handleCard(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: card <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  create <title> <card_number> <card_holder> <expiry> <cvv> [meta]")
		fmt.Println("  update <id> <title> <card_number> <card_holder> <expiry> <cvv> [meta]")
		fmt.Println("  delete <id>")
		fmt.Println("  list")
		os.Exit(1)
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "create":
		handleCreateCard(ctx, cli, subargs)
	case "update":
		handleUpdateCard(ctx, cli, subargs)
	case "delete":
		handleDeleteCard(ctx, cli, subargs)
	case "list":
		handleListCards(ctx, cli)
	default:
		fmt.Printf("Unknown card command: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleCreateCard(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 5 {
		fmt.Println("Usage: card create <title> <card_number> <card_holder> <expiry> <cvv> [meta]")
		os.Exit(1)
	}

	title := args[0]
	cardNumber := args[1]
	cardHolder := args[2]
	expiry := args[3]
	cvv := args[4]
	meta := ""
	if len(args) > 5 {
		meta = args[5]
	}

	card, err := cli.CreateCard(ctx, title, cardNumber, cardHolder, expiry, cvv, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create card: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Card created with ID: %s\n", card.GetId())
}

func handleUpdateCard(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 6 {
		fmt.Println("Usage: card update <id> <title> <card_number> <card_holder> <expiry> <cvv> [meta]")
		os.Exit(1)
	}

	id := args[0]
	title := args[1]
	cardNumber := args[2]
	cardHolder := args[3]
	expiry := args[4]
	cvv := args[5]
	meta := ""
	if len(args) > 6 {
		meta = args[6]
	}

	card, err := cli.UpdateCard(ctx, id, title, cardNumber, cardHolder, expiry, cvv, meta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update card: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Card updated: %s\n", card.GetId())
}

func handleDeleteCard(ctx context.Context, cli *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: card delete <id>")
		os.Exit(1)
	}

	id := args[0]

	if err := cli.DeleteCard(ctx, id); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete card: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Card deleted")
}

func handleListCards(ctx context.Context, cli *client.Client) {
	cards, err := cli.ListCards(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list cards: %v\n", err)
		os.Exit(1)
	}

	if len(cards) == 0 {
		fmt.Println("No cards found")
		return
	}

	fmt.Println("Cards:")
	for _, card := range cards {
		var maskedNumber string
		cardNumber := card.GetCardNumber()
		if len(cardNumber) >= 4 {
			maskedNumber = "****" + cardNumber[len(cardNumber)-4:]
		} else {
			maskedNumber = cardNumber
		}
		fmt.Printf("  ID: %s, Title: %s, Card Number: %s\n", card.GetId(), card.GetTitle(), maskedNumber)
	}
}

func handleSync(ctx context.Context, cli *client.Client, args []string) {
	resp, err := cli.Sync(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Sync failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sync completed at: %s\n", resp.GetServerTime().AsTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("Credentials: %d\n", len(resp.GetCredentials()))
	fmt.Printf("Text Data: %d\n", len(resp.GetTextData()))
	fmt.Printf("Binary Data: %d\n", len(resp.GetBinaryData()))
	fmt.Printf("Cards: %d\n", len(resp.GetCards()))
}

func printVersion() {
	fmt.Println("GophKeeper Client")
	fmt.Printf("Version: %s\n", version)

	buildDate := "unknown"
	commit := "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && commit == "unknown" {
				commit = setting.Value
			}
			if setting.Key == "vcs.time" && buildDate == "unknown" {
				buildDate = setting.Value
			}
		}
	}

	fmt.Printf("Build Date: %s\n", buildDate)
	fmt.Printf("Commit: %s\n", commit)
}

func printUsage() {
	fmt.Println("GophKeeper Password Manager CLI")
	fmt.Println()
	fmt.Println("Usage: gophkeeper <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  register <email>                Register a new user (password will be prompted)")
	fmt.Println("  login <email>                   Login to the server (password will be prompted)")
	fmt.Println("  credential <command> [args]     Manage credentials")
	fmt.Println("  text <command> [args]           Manage text data")
	fmt.Println("  binary <command> [args]         Manage binary data")
	fmt.Println("  card <command> [args]           Manage cards")
	fmt.Println("  sync                            Synchronize data")
	fmt.Println("  version                         Show version information")
	fmt.Println()
	fmt.Println("Global Options:")
	fmt.Println("  -server <address>               Server address (default: localhost:50051)")
	fmt.Println("  -tls                            Enable TLS")
	fmt.Println("  -tls-cert <path>                TLS certificate file path")
	fmt.Println("  -username <email>               Username/email (optional, can be provided as argument)")
	fmt.Println()
	fmt.Println("Note: Passwords are never passed as command-line arguments for security reasons.")
	fmt.Println("      They will be securely prompted during login/register operations.")
}
