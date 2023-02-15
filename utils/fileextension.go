package utils

import (
	"path/filepath"
	"strings"
)

// dot is required
func IsAllowedExtension(path string) bool {
	allowedExtensions := []string{
		".png",
		".jpg",
		".jpeg",
		".mov",
		".3gp",
		".mp4",
		".avi",
		".mkv",
	}

	return containsExtension(path, allowedExtensions)
}

// dot is required
func IsExifImage(path string) bool {
	allowedExtensions := []string{
		".jpg",
		".jpeg",
	}

	return containsExtension(path, allowedExtensions)
}

func GetExtension(path string) string {
	extension := filepath.Ext(path)

	return strings.ToLower(extension)
}

func containsExtension(path string, allowedExtensions []string) bool {
	extension := GetExtension(path)

	return contains(allowedExtensions, extension)
}

func contains(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}

	return false
}
