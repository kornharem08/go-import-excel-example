package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetLatestBackupFile returns the path to the latest backup file for the given source file
func GetLatestBackupFile(sourcePath string) (string, error) {
	backupDir := "backup"
	fileName := filepath.Base(sourcePath)
	backupPath := filepath.Join(backupDir, fileName)

	// Check if backup file exists
	_, err := os.Stat(backupPath)
	if err != nil {
		return "", fmt.Errorf("no backup file found: %v", err)
	}

	return backupPath, nil
}

// BackupFile copies a file to a backup folder, replacing any existing file with the same name
func BackupFile(sourcePath string) (string, error) {
	// Create backup directory if it doesn't exist
	backupDir := "backup"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Use the original filename for backup
	fileName := filepath.Base(sourcePath)
	backupPath := filepath.Join(backupDir, fileName)

	// Read source file
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %v", err)
	}

	// Write to backup file (will replace if exists)
	if err := os.WriteFile(backupPath, sourceData, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %v", err)
	}

	return backupPath, nil
}
