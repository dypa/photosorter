package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func FileContentHash(filePath string) (string, error) {
	hash := ""

	file, err := os.Open(filePath)
	if err != nil {
		return hash, err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return hash, err
	}
	hashInBytes := hasher.Sum(nil)
	hash = hex.EncodeToString(hashInBytes[:16])

	return hash, nil
}
