package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func CalculateChecksum(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:]), nil
}
