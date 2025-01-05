package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type Up struct {
	Post struct {
		File     *multipart.FileHeader `form:"file"`
		Name     string                `form:"name"`
		FileType string                `form:"fileType"`
		Path     string                `form:"path"`
	} `method:"POST"`
}

func (m *Up) HandleOptions(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

type upload struct {
	Code    int
	Data    map[string]interface{}
	Message string
}

func (m *Up) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	contextValue := contexext.FromContext(context)

	ossHost, err := context.Etcd().SelectInsideServer(config.MicroServerOSS)
	if err != nil {
		return nil, err
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", m.Post.File.Filename)
	if err != nil {
		return nil, err
	}
	fh, err := m.Post.File.Open()
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	dotIndex := strings.LastIndex(m.Post.File.Filename, ".")
	m.Post.FileType = m.Post.File.Filename[dotIndex:]
	err = bodyWriter.WriteField("name", m.Post.Name)
	if err != nil {
		return nil, err
	}
	log.Println(m.Post.File.Filename)

	err = bodyWriter.WriteField("fileType", m.Post.FileType)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%s", contextValue.DomainName, m.Post.Path)
	err = bodyWriter.WriteField("path", path)
	if err != nil {
		return nil, err
	}
	bodyWriter.Close()

	var post *http.Response
	post, err = http.Post(fmt.Sprintf("http://%s/upload", ossHost), bodyWriter.FormDataContentType(), bodyBuf)
	if err != nil {
		return nil, err
	}
	defer post.Body.Close()

	var u upload
	body, err := ioutil.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: u}, err
}

func (m *Up) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Up) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}
