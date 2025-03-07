package files

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, os.ModePerm)
}

func CopyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return err
			}
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			dstDir := filepath.Dir(dstPath)
			if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
				return err
			}
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func FileAlreadyExist(filePath string) bool {
	info := getEntryInfo(filePath)
	return info != nil
}

func DirAlreadyExist(folderPath string) bool {
	info := getEntryInfo(folderPath)
	if info == nil {
		return false
	}
	return info.IsDir()
}

func UpdateFile(filePath, searchText, replacementText string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	updatedContent := strings.ReplaceAll(string(data), searchText, replacementText)
	if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
		return err
	}
	return nil
}

func ContentExistsInFile(filePath, substringToFind string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(substringToFind)) {
			return true, nil // Found the substring, return true
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func CleanDestinationDir(dst string) error {
	return os.RemoveAll(dst)
}

func getEntryInfo(path string) os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		// An error occurred (e.g., permission denied)
		return nil
	}
	return info
}
