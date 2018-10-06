package service

import (
	"io/ioutil"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/nbvghost/gweb/tool"
)

type FileService struct{}

func (self FileService) DownNetImage(url string) string {

	resp, err := http.Get(url)
	tool.CheckError(err)
	if err != nil {
		return ""
	} else {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	tool.CheckError(err)
	return tool.WriteFile(b, resp.Header.Get("Content-Type"))
}
func (self FileService) DownNetWriteAliyunOSS(url string) string {

	resp, err := http.Get(url)
	tool.CheckError(err)
	if err != nil {
		return "no_found"
	}
	if resp.StatusCode == 404 {
		return "no_found"
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	tool.CheckError(err)

	ContentType := resp.Header.Get("Content-Type")

	path := tool.WriteTempFile(b, ContentType)
	//fmt.Println(path)
	if true {
		//return path
	}

	client, err := oss.New("oss-cn-shenzhen.aliyuncs.com", "tsZrY5eCZh9hQRbj", "CI3p9tiZGYHcN1wYgQBfZqqsAk6r8K")
	if err != nil {
		// HandleError(err)
		tool.CheckError(err)
	}

	bucket, err := client.Bucket("0e99ac3738124974b3ec74caf14f06fe")
	if err != nil {
		// HandleError(err)
		tool.CheckError(err)
	}

	err = bucket.PutObjectFromFile(path, "temp/"+path)
	if err != nil {
		// HandleError(err)
		tool.CheckError(err)
	}

	return "https://files.nutsy.cc/" + path
}
