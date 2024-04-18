package oss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"
)

type Upload struct {
	Code int
	Data struct {
		Path string
	}
	Message string
}

func Url(context constrain.IContext) (string, error) {
	ossHost, err := config.MicroServerOSS.SelectOutsideServer() //context.SelectOutsideServer(config.MicroServerOSS)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("//%s/assets", ossHost), nil
}
func ReadUrl(context constrain.IContext, path string) (string, error) {
	ossHost, err := config.MicroServerOSS.SelectOutsideServer() //context.SelectOutsideServer(config.MicroServerOSS)
	if err != nil {
		return "", err
	}
	contextValue := contexext.FromContext(context)
	return fmt.Sprintf("%s://%s/assets%s", util.GetScheme(contextValue.Request), ossHost, path), nil
}
func WriteUrl(context constrain.IContext) (string, error) {
	ossHost, err := config.MicroServerOSS.SelectInsideServer() //context.SelectInsideServer(config.MicroServerOSS)
	if err != nil {
		return "", err
	}
	//contextValue := contexext.FromContext(context)

	//内部连接，只有http,没有https
	return fmt.Sprintf("%s://%s/upload", "http", ossHost), nil
}

// UploadFile
// file 文件内容
//
// name 文件名，如果为空的话，使用md5(file)+fileType，做为文件名，如果fileType为空的话，则创建一个没有文件扩展名的文件
//
// fileType 只有在name 为空的情况下，才会生效
//
// override  如果存在相同的文件名时，是否覆盖原来的文件
// path 文件要存储的路径，做为name的目录
func UploadFile(context constrain.IContext, file []byte, path string, fileType string, override bool, name string) (*Upload, error) {
	ossUrl, err := WriteUrl(context)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = writer.WriteField("path", path)
	if err != nil {
		return nil, err
	}
	err = writer.WriteField("fileType", fileType)
	if err != nil {
		return nil, err
	}
	err = writer.WriteField("override", fmt.Sprintf("%t", override))
	if err != nil {
		return nil, err
	}
	err = writer.WriteField("name", name)
	if err != nil {
		return nil, err
	}
	formFile, err := writer.CreateFormFile("file", "filename")
	if err != nil {
		return nil, err
	}
	_, err = formFile.Write(file)
	if err != nil {
		return nil, err
	}

	//必须执行，否则http不知道是否内容结束
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	post, err := http.Post(ossUrl, writer.FormDataContentType(), buf)
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()

	body, err := ioutil.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}
	var upload Upload
	err = json.Unmarshal(body, &upload)
	if err != nil {
		return nil, err
	}
	return &upload, nil
}
func UploadAvatar(context constrain.IContext, OID, userID dao.PrimaryKey, file []byte) (*Upload, error) {

	return UploadFile(context, file, fmt.Sprintf("miniapp/avatar/%d/%d", OID, userID), "", false, fmt.Sprintf("%s", time.Now().Format("20060102150405")))

}
