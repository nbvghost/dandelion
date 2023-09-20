package file

import (
	"github.com/nbvghost/tool/encryption"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type FileService struct{}

func (m FileService) DownNetImage(url string) (string, error) {

	resp, err := http.Get(url)
	log.Println(err)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	log.Println(err)
	return m.WriteFile(b, resp.Header.Get("Content-Type"), "", "")
}
func (m FileService) WriteFile(fileBytes []byte, ContentType string, dynamicDirName, dirType string) (string, error) {
	md5Name := encryption.Md5ByBytes(fileBytes)

	var f *os.File

	fileTypes := strings.Split(ContentType, "/")
	if len(fileTypes) != 2 {

		return "", errors.New("ContentType 格式不正确")
	}

	contentTypeList := strings.Split(fileTypes[1], "-")
	fileType := contentTypeList[len(contentTypeList)-1]
	filePath := strings.Trim(dynamicDirName, "/") + "/" + strings.Trim(dirType, "/") + "/"
	fileFullPath := strings.TrimRight("upload", "/") + "/" + filePath
	if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
		os.MkdirAll(fileFullPath, os.ModePerm)
	}

	fileName := fileFullPath + "/" + md5Name + "." + fileType

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		//不存在的文件
		f, err = os.Create(fileName) //创建文件
		if err != nil {
			return "", err
		}
		defer f.Close()
		f.Write(fileBytes)
		f.Sync()

	}
	return filePath + "/" + md5Name + "." + fileType, nil
}
func (m FileService) DownNetWriteAliyunOSS(url string) (string, error) {

	resp, err := http.Get(url)
	log.Println(err)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 404 {
		return "", errors.New("no_found")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ContentType := resp.Header.Get("Content-Type")

	path, err := m.WriteTempFile(b, ContentType)
	if err != nil {
		return "", err
	}
	//fmt.Println(path)
	if true {
		//return path
	}

	client, err := oss.New("oss-cn-shenzhen.aliyuncs.com", "tsZrY5eCZh9hQRbj", "CI3p9tiZGYHcN1wYgQBfZqqsAk6r8K")
	if err != nil {
		// HandleError(err)
		return "", err
	}

	bucket, err := client.Bucket("0e99ac3738124974b3ec74caf14f06fe")
	if err != nil {
		// HandleError(err)
		return "", err
	}

	err = bucket.PutObjectFromFile(path, "temp/"+path)
	if err != nil {
		// HandleError(err)
		return "", err
	}

	return "https://files.nutsy.cc/" + path, nil
}
func (m FileService) WriteTempFile(fileBytes []byte, ContentType string) (string, error) {

	md5Name := strings.ToLower(encryption.Md5ByBytes(fileBytes))
	var f *os.File

	fileTypes := strings.Split(ContentType, "/")
	if len(fileTypes) != 2 {
		return "", errors.New("ContentType 格式不正确")
	}
	contentTypeList := strings.Split(fileTypes[1], "-")
	fileType := contentTypeList[len(contentTypeList)-1]

	fileName := md5Name + "." + fileType
	fullPath := "temp/" + fileName

	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.MkdirAll("temp", os.ModePerm)
	}

	f, err := os.Create(fullPath) //创建文件
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.Write(fileBytes)
	f.Sync()

	return fileName, nil

}
