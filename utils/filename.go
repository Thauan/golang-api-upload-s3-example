package utils

import "path/filepath"

func GetFileWithDir(key string) (string, string) {
	dir, file := filepath.Split(key)
	return dir, file
}
