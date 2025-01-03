package controller

import (
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

type OSSUpload struct {
	Get struct {
		//TempFilename string `form:"filename"`
		Path string `form:"path"`
	} `method:"Get"`
}

func (m *OSSUpload) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)
	err := contextValue.Request.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		return nil, err
	}

	var file multipart.File
	file, fileHeader, err := contextValue.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileByte, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(fileHeader.Filename)
	/*if strings.EqualFold(filepath.Ext(fileHeader.Filename), ".mp4") {

	} else {

	}
	buffer := bytes.NewBuffer(nil)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.Errorf("不支持的图片格式，请确认图片是否合格,%s", err.Error())
	}
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return nil, errors.Errorf("不支持的图片格式，请确认图片是否合格,%s", err.Error())
	}*/

	//fileByte := buffer.Bytes()
	fileName, err := oss.CreateTempWithExt(fileByte, strings.ToLower(ext))
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Path": fileName, "Name": fileHeader.Filename}), nil
}

func (m *OSSUpload) Handle(context constrain.IContext) (constrain.IResult, error) {
	if strings.HasPrefix(m.Get.Path, oss.TempFilePrefix) {
		b, err := oss.GetTempFile(m.Get.Path)
		if err != nil {
			return nil, err
		}
		return &result.ImageBytesResult{
			Data:        b,
			ContentType: "",
			Filename:    m.Get.Path,
		}, nil
	} else {
		return m.ossLoad(context)
	}

}

func (m *OSSUpload) ossLoad(ctx constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(ctx)
	server, err := ctx.Etcd().SelectInsideServer(config.MicroServerOSS) //config.MicroServerOSS.SelectInsideServer() //ctx.SelectInsideServer(config.MicroServerOSS)
	if err != nil {
		return nil, err
	}
	ossUrl, err := url.Parse(fmt.Sprintf("http://%s", server))
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(ossUrl)
	reverseProxy.Director = func(request *http.Request) {

		targetQuery := ossUrl.RawQuery
		director := func(req *http.Request) {
			req.URL.Scheme = ossUrl.Scheme
			req.URL.Host = ossUrl.Host
			req.URL.Path, req.URL.RawPath = joinURLPath(ossUrl, req.URL)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
			path := filepath.ToSlash(filepath.Join("/assets", m.Get.Path))
			req.URL.Path = path
		}
		director(request)
	}
	reverseProxy.ServeHTTP(contextValue.Response, contextValue.Request)

	return &result.NoneResult{}, nil
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
