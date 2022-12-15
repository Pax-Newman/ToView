package filehelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetFilename(path string) string {
	base := filepath.Base(path)
	name := strings.Split(base, ".")[0]
	return name
}

func GetExtension(path string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return "", fmt.Errorf("file \"%s\" does not have an extension and cannot be parsed", path)
	}
	return ext[1:], nil
}

// Check if a file is valid, supported, and exists
func CheckValid(path string) error {
	// check if the filepath is valid
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	// check if the file extension is valid
	_, err = GetExtension(path)
	if err != nil {
		return err
	}

	return nil
}
