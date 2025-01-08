package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func hashFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func generateManifest(dirPath string, manifestPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	manifest := make(map[string]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(dirPath, file.Name())
		hash, err := hashFileSHA256(filePath)
		if err != nil {
			return fmt.Errorf("failed to hash file %s: %w", file.Name(), err)
		}
		manifest[file.Name()] = hash
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest to JSON: %w", err)
	}

	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file %s: %w", manifestPath, err)
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <directory path> <manifest output path>\n", os.Args[0])
	}

	dirPath := os.Args[1]
	manifestPath := os.Args[2]

	err := generateManifest(dirPath, manifestPath)
	if err != nil {
		log.Fatalf("Error generating manifest: %v\n", err)
	}

	fmt.Println("Manifest file created successfully.")
}
