package file

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nbvghost/gweb"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type FileService struct{}

func (self FileService) DownNetImage(url string) string {

	resp, err := http.Get(url)
	log.Println(err)
	if err != nil {
		return ""
	} else {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	log.Println(err)
	return gweb.WriteFile(b, resp.Header.Get("Content-Type"), "", "")
}
func (self FileService) DownNetWriteAliyunOSS(url string) string {

	resp, err := http.Get(url)
	log.Println(err)
	if err != nil {
		return "no_found"
	}
	if resp.StatusCode == 404 {
		return "no_found"
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	log.Println(err)

	ContentType := resp.Header.Get("Content-Type")

	path := gweb.WriteTempFile(b, ContentType)
	//fmt.Println(path)
	if true {
		//return path
	}

	client, err := oss.New("oss-cn-shenzhen.aliyuncs.com", "tsZrY5eCZh9hQRbj", "CI3p9tiZGYHcN1wYgQBfZqqsAk6r8K")
	if err != nil {
		// HandleError(err)
		log.Println(err)
	}

	bucket, err := client.Bucket("0e99ac3738124974b3ec74caf14f06fe")
	if err != nil {
		// HandleError(err)
		log.Println(err)
	}

	err = bucket.PutObjectFromFile(path, "temp/"+path)
	if err != nil {
		// HandleError(err)
		log.Println(err)
	}

	return "https://files.nutsy.cc/" + path
}
