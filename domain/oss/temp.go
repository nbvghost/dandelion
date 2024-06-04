package oss

import (
	"github.com/nbvghost/tool/encryption"
	"os"
	"strings"
)

func CreateTempFilename(fileByte []byte) (string, error) {
	fileMD5 := strings.ToLower(encryption.Md5ByBytes(fileByte))
	filePath := os.TempDir() + "/" + fileMD5
	err := os.WriteFile(filePath, fileByte, os.ModePerm)
	if err != nil {
		return "", err
	}
	return fileMD5, nil
}
func GetTempFile(filename string) ([]byte, error) {
	filePath := os.TempDir() + "/" + filename
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}