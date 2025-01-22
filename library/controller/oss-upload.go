package controller

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

type OSSUpload struct {
	Admin *entity.SessionMappingData `mapping:""`
	Get   struct {
		Path string `form:"path"`
	} `method:"Get"`
	Post struct {
		Direct   bool           `form:"Direct"`
		TargetID dao.PrimaryKey `form:"TargetID"`
		Target   string         `form:"Target"`
		Title    string         `form:"Title"`
	} `method:"Post"`
}

func (m *OSSUpload) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)

	if m.Post.Direct {
		if m.Post.TargetID == 0 {
			return nil, errors.New("TargetID不能为空")
		}
		if len(m.Post.Target) == 0 {
			return nil, errors.New("Target不能为空")
		}
		uploadFile, err := oss.UploadFileProxy(context, contextValue.Response, contextValue.Request)
		if err != nil {
			return nil, err
		}
		if uploadFile.Code != 0 {
			return nil, result.NewCodeWithMessage(result.ActionResultCode(uploadFile.Code), uploadFile.Message)
		}
		media := dao.GetBy(db.Orm(), &model.Media{}, map[string]any{"OID": m.Admin.OID, "TargetID": m.Post.TargetID, "Target": m.Post.Target, "SHA256": uploadFile.Data.SHA256}).(*model.Media)
		if media.IsZero() {
			media = &model.Media{
				OID:      m.Admin.OID,
				TargetID: m.Post.TargetID,
				Target:   model.MediaTarget(m.Post.Target),
				SHA256:   uploadFile.Data.SHA256,
				Title:    m.Post.Title,
				Src:      uploadFile.Data.Path,
				Size:     uploadFile.Data.Size,
				Width:    uploadFile.Data.Width,
				Height:   uploadFile.Data.Height,
				FileName: uploadFile.Data.Filename,
				Format:   uploadFile.Data.Format,
				Tags:     []string{},
			}
			err = dao.Create(db.Orm(), media)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	} else {

		err := contextValue.Request.ParseMultipartForm(10 * 1024 * 1024)
		if err != nil {
			return nil, err
		}
		var file multipart.File
		var fileHeader *multipart.FileHeader
		file, fileHeader, err = contextValue.Request.FormFile("file")
		if err != nil {
			return nil, err
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)

		fileByte, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		ext := filepath.Ext(fileHeader.Filename)

		fileName, err := oss.CreateTempWithExt(fileByte, strings.ToLower(ext))
		if err != nil {
			return nil, err
		}
		return result.NewData(map[string]any{"Path": fileName, "Name": fileHeader.Filename}), nil
	}
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
	/*contextValue := contexext.FromContext(ctx)
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
	reverseProxy.ServeHTTP(contextValue.Response, contextValue.Request)*/

	contextValue := contexext.FromContext(ctx)
	server, err := ctx.Etcd().SelectInsideServer(config.MicroServerOSS) //config.MicroServerOSS.SelectInsideServer() //ctx.SelectInsideServer(config.MicroServerOSS)
	if err != nil {
		return nil, err
	}
	ossUrl, err := url.Parse(fmt.Sprintf("http://%s/assets%s", server, m.Get.Path))
	if err != nil {
		return nil, err
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(ossUrl)
	reverseProxy.Director = func(req *http.Request) {
		req.URL = ossUrl
	}
	reverseProxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Del("Access-Control-Allow-Origin")
		response.Header.Del("Access-Control-Allow-Headers")
		response.Header.Del("Access-Control-Allow-Methods")
		response.Header.Del("Access-Control-Allow-Credentials")
		return nil
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
