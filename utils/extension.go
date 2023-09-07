package utils

import "path/filepath"

func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return ext
}

// Função para verificar se um valor está em uma lista de valores
func ContainsFileExtension(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
