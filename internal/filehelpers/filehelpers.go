package filehelpers

import (
	"fmt"
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
