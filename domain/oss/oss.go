package oss

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/tool/encryption"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/util"
)

type Upload struct {
	Code int
	Data struct {
		Path     string
		Ext      string
		Filename string
		Size     int
		Width    int
		Height   int
		Format   string
		SHA256   string
	}
	Message string
}

func Url(context constrain.IContext) (string, error) {
	ossHost, err := context.Etcd().SelectOutsideServer(config.MicroServerOSS)
	if err != nil {
		return "", err
	}
	contextValue := contexext.FromContext(context)
	return fmt.Sprintf("%s://%s/assets", util.GetScheme(contextValue.Request), ossHost), nil
}
func ReadUrl(context constrain.IContext, path string) (string, error) {
	ossHost, err := context.Etcd().SelectOutsideServer(config.MicroServerOSS)
	if err != nil {
		return "", err
	}
	contextValue := contexext.FromContext(context)
	return fmt.Sprintf("%s://%s/assets%s", util.GetScheme(contextValue.Request), ossHost, path), nil
}
func WriteUrl(context constrain.IServiceContext) (string, error) {
	ossHost, err := context.Etcd().SelectInsideServer(config.MicroServerOSS)
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
func UploadFileBase(ossUrl string, file []byte, path string, fileType string, override bool, name string) (*Upload, error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err := writer.WriteField("path", path)
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
	filename := name
	if len(filename) == 0 {
		filename = strings.ToLower(encryption.Md5ByBytes(file))
	}
	formFile, err := writer.CreateFormFile("file", filename)
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
func UploadFile(context constrain.IServiceContext, file []byte, path string, fileType string, override bool, name string) (*Upload, error) {
	ossUrl, err := WriteUrl(context)
	if err != nil {
		return nil, err
	}
	return UploadFileBase(ossUrl, file, path, fileType, override, name)
}

type ProxyResponse struct {
}

func (m *ProxyResponse) Header() http.Header {
	return http.Header{}
}

func (m *ProxyResponse) Write(i []byte) (int, error) {

	return len(i), nil
}

func (m *ProxyResponse) WriteHeader(statusCode int) {
}

func UploadFileProxy(context constrain.IServiceContext, writer http.ResponseWriter, request *http.Request) (*Upload, error) {
	//contextValue := contexext.FromContext(context)
	server, err := context.Etcd().SelectInsideServer(config.MicroServerOSS) //config.MicroServerOSS.SelectInsideServer() //ctx.SelectInsideServer(config.MicroServerOSS)
	if err != nil {
		return nil, err
	}
	ossUrl, err := url.Parse(fmt.Sprintf("http://%s/upload", server))
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(ossUrl)
	reverseProxy.Director = func(req *http.Request) {
		req.URL = ossUrl
	}

	response := &ProxyResponse{}
	var callback = make(chan *Upload, 1)
	reverseProxy.ModifyResponse = func(response *http.Response) error {
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var upload Upload
		err = json.Unmarshal(body, &upload)
		if err != nil {
			return err
		}
		callback <- &upload
		//response.Header.Del("Access-Control-Allow-Origin")
		//response.Header.Del("Access-Control-Allow-Headers")
		//response.Header.Del("Access-Control-Allow-Methods")
		//response.Header.Del("Access-Control-Allow-Credentials")
		return nil
	}
	reverseProxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		callback <- &Upload{Code: 9077, Message: err.Error()}
	}
	reverseProxy.ServeHTTP(response, request)

	t := time.NewTicker(time.Second * 60)
	for {
		select {
		case upload := <-callback:
			return upload, nil
		case <-t.C:
			return nil, errors.New("oss upload file timeout")
		}
	}
}
func UploadAvatar(context constrain.IContext, OID, userID dao.PrimaryKey, file []byte) (*Upload, error) {

	return UploadFile(context, file, fmt.Sprintf("miniapp/avatar/%d/%d", OID, userID), "", false, fmt.Sprintf("%s", time.Now().Format("20060102150405")))

}
