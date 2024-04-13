// Copyright (c) 2024. Hangover Games <info@hangover.games>. All rights reserved.

package utils

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// Embedding the template Go file
//
//go:embed templates/main.go
var embeddedFile embed.FS

// CompileBinary compiles the provided encrypted binary data as a self-executable installer
func CompileBinary(encryptedData []byte, outputPath string) error {

	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "example")
	if err != nil {
		return fmt.Errorf("failed creating temporary directory: %w", err)
	}

	// Don't forget to clean up after you're done
	defer os.RemoveAll(tempDir)

	// Create a temporary Go file with embedded data
	// This would involve writing the encrypted data as a byte slice in Go source code
	tempExecutablePath := filepath.Join(tempDir, "encrypted.txt")

	// Create a temporary Go file with embedded data
	// This would involve writing the encrypted data as a byte slice in Go source code
	tempFilePath := filepath.Join(tempDir, "main.go")

	// Load the template source
	sourceCode, err := fs.ReadFile(embeddedFile, "templates/main.go")
	if err != nil {
		return fmt.Errorf("failed to read embedded source file: %w", err)
	}

	// Create embedded encrypted.txt
	err = os.WriteFile(tempExecutablePath, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to the encrypted file '%s': %w", tempExecutablePath, err)
	}

	// Create embedded source file
	err = os.WriteFile(tempFilePath, sourceCode, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to the source file '%s': %w", tempFilePath, err)
	}

	// Compile the Go program
	var stderr bytes.Buffer
	cmd := exec.Command("go", "build", "-ldflags", "-extldflags \"-static\"", "-o", outputPath, tempFilePath)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to compile: %w; stderr:\n%s", err, stderr.String())
	}

	return nil
}
