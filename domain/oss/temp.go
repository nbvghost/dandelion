package oss

import (
	"github.com/nbvghost/tool/encryption"
	"os"
	"strings"
)

const TempFilePrefix string = "temp-file:"

func CreateTempWithExt(fileByte []byte, ext string) (string, error) {
	fileName := TempFilePrefix + strings.ToLower(encryption.Md5ByBytes(fileByte)) + ext
	filePath := os.TempDir() + "/" + fileName
	err := os.WriteFile(filePath, fileByte, os.ModePerm)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func CreateTempFilename(fileByte []byte) (string, error) {
	return CreateTempWithExt(fileByte, "")
	//fileMD5 := strings.ToLower(encryption.Md5ByBytes(fileByte))
	//filePath := os.TempDir() + "/" + fileMD5
	//err := os.WriteFile(filePath, fileByte, os.ModePerm)
	//if err != nil {
	//	return "", err
	//}
	//return fileMD5, nil
}
func GetTempFile(filename string) ([]byte, error) {
	filePath := os.TempDir() + "/" + filename
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}
