package utils

import "path/filepath"

func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return ext
}
