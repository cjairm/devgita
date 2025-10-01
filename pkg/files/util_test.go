package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/logger"
	"github.com/cjairm/devgita/pkg/files"
)

var fileContent = "Hello, World!"

func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "testdir-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	return tempDir
}

func createTempSourceFile(t *testing.T, tempDir string, fileName string) string {
	srcFilePath := filepath.Join(tempDir, fileName)
	err := os.WriteFile(srcFilePath, []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	return srcFilePath
}

func createTempSourceDir(t *testing.T, tempDir string) string {
	srcDir := filepath.Join(tempDir, "source")
	err := os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	createTempSourceFile(t, srcDir, "file1.txt")

	subDir := filepath.Join(srcDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	createTempSourceFile(t, subDir, "file2.txt")

	return srcDir
}

func TestCopyFile(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)
	logger.Init(false)

	srcFilePath := createTempSourceFile(t, tempDir, "source.txt")
	dstFilePath := filepath.Join(tempDir, "destination.txt")

	// Test successful copy
	if err := files.CopyFile(srcFilePath, dstFilePath); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}
	t.Log("File copied successfully")

	// Verify the content of the destination file
	dstContent, err := os.ReadFile(dstFilePath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	expectedContent := "Hello, World!"
	if string(dstContent) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, dstContent)
	}
	t.Log("Content copied succesfuly")

	// Test copying from a non-existent source file
	nonExistentSrc := filepath.Join(tempDir, "nonexistent.txt")
	err = files.CopyFile(nonExistentSrc, dstFilePath)
	if err == nil {
		t.Fatal("Expected an error when copying from a nonexistent file, got nil")
	}
	t.Log("Correctly handled copying from a nonexistent file")

	// Test copying to a directory instead of a file
	dstDirPath := filepath.Join(tempDir, "destinationDir")
	if err := os.Mkdir(dstDirPath, 0755); err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}
	err = files.CopyFile(srcFilePath, dstDirPath)
	if err == nil {
		t.Fatal("Expected an error when copying to a directory, got nil")
	}
	t.Log("Correctly handled copying to a directory")
}

func TestCopyDir(t *testing.T) {
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a source directory with files
	srcDir := createTempSourceDir(t, tempDir)
	dstDir := filepath.Join(tempDir, "destination")

	// Test successful copy
	if err := files.CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}
	t.Log("Successfully copied files from source to destination.")

	// Verify the content of the destination directory
	dstFile1 := filepath.Join(dstDir, "file1.txt")
	dstFile2 := filepath.Join(dstDir, "subdir", "file2.txt")

	if _, err := os.Stat(dstFile1); os.IsNotExist(err) {
		t.Errorf("Expected file %q to exist, but it does not", dstFile1)
	}

	if _, err := os.Stat(dstFile2); os.IsNotExist(err) {
		t.Errorf("Expected file %q to exist, but it does not", dstFile2)
	}

	// Test copying an empty directory
	emptySrcDir := filepath.Join(srcDir, "empty")
	if err := os.Mkdir(emptySrcDir, 0755); err != nil {
		t.Fatalf("Failed to create empty source directory: %v", err)
	}
	emptyDstDir := filepath.Join(dstDir, "empty")
	if err := files.CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDir failed for empty directory: %v", err)
	}
	if _, err := os.Stat(emptyDstDir); os.IsNotExist(err) {
		t.Errorf("Expected empty destination directory %q to exist, but it does not", emptyDstDir)
	}

	// Test copying from a non-existent source directory
	nonExistentSrcDir := filepath.Join(tempDir, "nonexistent")
	err := files.CopyDir(nonExistentSrcDir, dstDir)
	if err == nil {
		t.Fatal("Expected an error when copying from a nonexistent directory, got nil")
	}
	t.Log("Correctly received an error when copying from a nonexistent directory.")

	// Test copying to an existing directory
	existingFile := filepath.Join(dstDir, "file1.txt")
	if err := os.WriteFile(existingFile, []byte("Existing file"), 0644); err != nil {
		t.Fatalf("Failed to create existing file in destination: %v", err)
	}
	if err := files.CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDir failed when copying to an existing directory: %v", err)
	}
	t.Log("Successfully copied files to an existing directory.")

	// Verify that the existing file is still there
	if _, err := os.Stat(existingFile); os.IsNotExist(err) {
		t.Errorf("Expected existing file %q to still exist, but it does not", existingFile)
	}

	// Test handling read-only files
	readOnlyFile := filepath.Join(dstDir, "readonly.txt")
	if err := os.WriteFile(readOnlyFile, []byte("Read-only file"), 0444); err != nil {
		t.Fatalf("Failed to create read-only file: %v", err)
	}
	if err := files.CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDir failed when copying with read-only file: %v", err)
	}
	t.Log("Successfully copied files with a read-only file present.")

	// Verify that the read-only file is still there
	if _, err := os.Stat(readOnlyFile); os.IsNotExist(err) {
		t.Errorf("Expected read-only file %q to still exist, but it does not", readOnlyFile)
	}
}
