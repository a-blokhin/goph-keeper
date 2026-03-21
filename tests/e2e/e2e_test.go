package e2e

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	testEmail    = "e2e-test@example.com"
	testPassword = "testpassword123"
)

var (
	clientPath string
	serverAddr = "localhost:50051"
)

func TestMain(m *testing.M) {
	if os.Getenv("RUN_E2E") != "1" {
		os.Exit(0)
	}
	var err error

	clientPath, err = filepath.Abs("../../bin/client")
	if err != nil {
		fmt.Printf("Failed to get client path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		fmt.Printf("Client binary not found at %s. Please build it first: go build -o bin/client ./cmd/client\n", clientPath)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func runClientCommand(ctx context.Context, args ...string) (string, error) {
	fullArgs := []string{"--server", serverAddr}
	fullArgs = append(fullArgs, args...)

	cmd := exec.CommandContext(ctx, clientPath, fullArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Env = os.Environ()

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func TestE2EFullWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	t.Run("RegisterUser", func(t *testing.T) {
		output, err := runClientCommand(ctx, "register", testEmail, testPassword)
		if err != nil {
			t.Logf("Registration failed (user may already exist), skipping: %v", err)
		} else if !strings.Contains(output, "Registration successful") {
			t.Fatalf("Unexpected output: %s", output)
		}
	})

	t.Run("Login", func(t *testing.T) {
		output, err := runClientCommand(ctx, "login", testEmail, testPassword)
		if err != nil {
			t.Fatalf("Failed to login: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Login successful") {
			t.Fatalf("Unexpected output: %s", output)
		}
	})

	var credentialID, textID, binaryID, cardID string

	t.Run("CreateCredential", func(t *testing.T) {
		output, err := runClientCommand(ctx, "credential", "create", "Test Website", "testuser", "testpass", "test.com")
		if err != nil {
			t.Fatalf("Failed to create credential: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Credential created") {
			t.Fatalf("Unexpected output: %s", output)
		}

		parts := strings.Fields(output)
		for i, part := range parts {
			if part == "ID:" && i+1 < len(parts) {
				credentialID = parts[i+1]
				break
			}
		}
		if credentialID == "" {
			t.Fatal("Could not extract credential ID")
		}
	})

	t.Run("ListCredentials", func(t *testing.T) {
		output, err := runClientCommand(ctx, "credential", "list")
		if err != nil {
			t.Fatalf("Failed to list credentials: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Test Website") {
			t.Fatalf("Credential not found in list: %s", output)
		}
	})

	t.Run("UpdateCredential", func(t *testing.T) {
		output, err := runClientCommand(ctx, "credential", "update", credentialID, "Updated Website", "newuser", "newpass", "updated.com")
		if err != nil {
			t.Fatalf("Failed to update credential: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Credential updated") {
			t.Fatalf("Unexpected output: %s", output)
		}
	})

	t.Run("CreateTextData", func(t *testing.T) {
		output, err := runClientCommand(ctx, "text", "create", "Test Note", "This is a test note", "personal")
		if err != nil {
			t.Fatalf("Failed to create text data: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Text data created") {
			t.Fatalf("Unexpected output: %s", output)
		}

		parts := strings.Fields(output)
		for i, part := range parts {
			if part == "ID:" && i+1 < len(parts) {
				textID = parts[i+1]
				break
			}
		}
		if textID == "" {
			t.Fatal("Could not extract text data ID")
		}
	})

	t.Run("ListTextData", func(t *testing.T) {
		output, err := runClientCommand(ctx, "text", "list")
		if err != nil {
			t.Fatalf("Failed to list text data: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Test Note") {
			t.Fatalf("Text data not found in list: %s", output)
		}
	})

	t.Run("UpdateTextData", func(t *testing.T) {
		output, err := runClientCommand(ctx, "text", "update", textID, "Updated Note", "Updated content", "updated")
		if err != nil {
			t.Fatalf("Failed to update text data: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Text data updated") {
			t.Fatalf("Unexpected output: %s", output)
		}
	})

	t.Run("CreateBinaryData", func(t *testing.T) {
		tmpFile := filepath.Join(os.TempDir(), "test_binary.bin")
		testData := []byte("test binary data")
		if err := os.WriteFile(tmpFile, testData, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer os.Remove(tmpFile)

		output, err := runClientCommand(ctx, "binary", "create", "Test File", tmpFile, "documents")
		if err != nil {
			t.Fatalf("Failed to create binary data: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Binary data created") {
			t.Fatalf("Unexpected output: %s", output)
		}

		parts := strings.Fields(output)
		for i, part := range parts {
			if part == "ID:" && i+1 < len(parts) {
				binaryID = parts[i+1]
				break
			}
		}
		if binaryID == "" {
			t.Fatal("Could not extract binary data ID")
		}
	})

	t.Run("ListBinaryData", func(t *testing.T) {
		output, err := runClientCommand(ctx, "binary", "list")
		if err != nil {
			t.Fatalf("Failed to list binary data: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Test File") {
			t.Fatalf("Binary data not found in list: %s", output)
		}
	})

	t.Run("CreateCard", func(t *testing.T) {
		output, err := runClientCommand(ctx, "card", "create", "Test Card", "1234567890123456", "John Doe", "12/25", "123", "personal")
		if err != nil {
			t.Fatalf("Failed to create card: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Card created") {
			t.Fatalf("Unexpected output: %s", output)
		}

		parts := strings.Fields(output)
		for i, part := range parts {
			if part == "ID:" && i+1 < len(parts) {
				cardID = parts[i+1]
				break
			}
		}
		if cardID == "" {
			t.Fatal("Could not extract card ID")
		}
	})

	t.Run("ListCards", func(t *testing.T) {
		output, err := runClientCommand(ctx, "card", "list")
		if err != nil {
			t.Fatalf("Failed to list cards: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Test Card") {
			t.Fatalf("Card not found in list: %s", output)
		}
	})

	t.Run("UpdateCard", func(t *testing.T) {
		output, err := runClientCommand(ctx, "card", "update", cardID, "Updated Card", "9876543210987654", "Jane Doe", "06/26", "456", "updated")
		if err != nil {
			t.Fatalf("Failed to update card: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Card updated") {
			t.Fatalf("Unexpected output: %s", output)
		}
	})

	t.Run("Sync", func(t *testing.T) {
		output, err := runClientCommand(ctx, "sync")
		if err != nil {
			t.Fatalf("Failed to sync: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "Sync completed") {
			t.Fatalf("Unexpected output: %s", output)
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		if credentialID != "" {
			_, err := runClientCommand(ctx, "credential", "delete", credentialID)
			if err != nil {
				t.Logf("Warning: Failed to delete credential: %v", err)
			}
		}

		if textID != "" {
			_, err := runClientCommand(ctx, "text", "delete", textID)
			if err != nil {
				t.Logf("Warning: Failed to delete text data: %v", err)
			}
		}

		if binaryID != "" {
			_, err := runClientCommand(ctx, "binary", "delete", binaryID)
			if err != nil {
				t.Logf("Warning: Failed to delete binary data: %v", err)
			}
		}

		if cardID != "" {
			_, err := runClientCommand(ctx, "card", "delete", cardID)
			if err != nil {
				t.Logf("Warning: Failed to delete card: %v", err)
			}
		}

		tokenFile := filepath.Join(os.Getenv("HOME"), ".gophkeeper_token")
		if err := os.Remove(tokenFile); err != nil && !os.IsNotExist(err) {
			t.Logf("Warning: Failed to remove token file: %v", err)
		}
	})
}

func TestE2EVersion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	output, err := runClientCommand(ctx, "version")
	if err != nil {
		t.Fatalf("Failed to get version: %v, output: %s", err, output)
	}
	if !strings.Contains(output, "Version:") {
		t.Fatalf("Unexpected output: %s", output)
	}
}
