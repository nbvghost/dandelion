package oss

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/domain/webpicture"
	"github.com/nbvghost/dandelion/server/httpext"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed default.png
var defaultImageBytes []byte

type TemplateDir struct {
}

func (t *TemplateDir) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	writer.Header().Set("Age", fmt.Sprintf("%d", 12*60*60))
	//public, max-age=31536000
	writer.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", 24*60*60))
	//log.Println("Origin:", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Width,Height")
	writer.Header().Set("Access-Control-Allow-Methods", "GET")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")

	request.URL.Path = strings.Split(request.URL.Path, "@")[0]

	fp := filepath.Join("assets", request.URL.Path)
	//log.Println("fp", fp)
	fileInfo, err := os.Stat(fp)
	if os.IsNotExist(err) || fileInfo.IsDir() {
		writeToDefaultImage(writer)
	} else {
		http.FileServer(http.Dir("assets")).ServeHTTP(writer, request)
	}
}
func writeToImage(writer http.ResponseWriter, body []byte, contentType string, w, h int) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Width", fmt.Sprintf("%d", w))
	writer.Header().Set("Height", fmt.Sprintf("%d", h))
	_, err := writer.Write(body)
	if err != nil {
		log.Println(err)
	}
}
func writeToErrImage(writer http.ResponseWriter, request *http.Request, err error) {
	http.NotFoundHandler().ServeHTTP(writer, request)
}
func writeToDefaultImage(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "image/png")
	writer.Header().Set("Width", fmt.Sprintf("%d", 1500))
	writer.Header().Set("Height", fmt.Sprintf("%d", 1500))
	_, err := writer.Write(defaultImageBytes)
	if err != nil {
		log.Println(err)
	}
}

type Server struct {
	etcdService constrain.IEtcd
	engine      *mux.Router
}

func (m *Server) Listen(ip string, port int) error {
	serverDesc := &config.MicroServerConfig{
		MicroServer: config.MicroServerOSS,
		IP:          ip,
		Port:        port,
	}
	httpServer := httpext.NewHttpServer(m.etcdService, nil, m.engine, nil, nil)
	//httpext.WithServerDesc(serverDesc.MicroServer.Name, serverDesc.IP, serverDesc.Port),
	log.Println("start oss server")
	return httpServer.Listen(serverDesc)
}

func NewServer(etcdService constrain.IEtcd) *Server {
	if err := os.Mkdir("assets", os.ModePerm); err != nil {
		var pathErr *fs.PathError
		ok := errors.As(err, &pathErr)
		if !ok {
			log.Fatalln(pathErr)
		}
	}

	_ = bmp.Decode
	_ = tiff.Decode
	_ = jpeg.Decode
	_ = png.Decode
	_ = gif.Decode
	_ = webp.Decode

	engine := mux.NewRouter()
	/*engine.HandleFunc("/browse", func(writer http.ResponseWriter, request *http.Request) {
		t, err := template.ParseFiles("browse.html")
		if err != nil {
			log.Println(err)
			return
		}
		err = t.Execute(writer, nil)
		if err != nil {
			log.Println(err)
			return
		}
	})*/
	engine.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", &TemplateDir{}))
	engine.Handle("/upload", &UploadHandle{})
	engine.Handle("/push/json", &PushJSON{})
	return &Server{etcdService: etcdService, engine: engine}
}

type PushJSON struct {
}

func (m *PushJSON) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	t := request.URL.Query().Get("t")
	body, _ := io.ReadAll(request.Body)
	fileName := fmt.Sprintf("%s-%s-%s.json", name, t, time.Now().Format(time.RFC3339))
	fileFullName := filepath.Join("assets", "push", "json", name, t, fileName)
	os.MkdirAll(filepath.Dir(fileFullName), os.ModePerm)
	os.WriteFile(fileFullName, body, os.ModePerm)
	writer.Write([]byte("SUCCESS"))
}

type UploadHandle struct {
}

func (m *UploadHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")

	err := request.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		//result.Code = 9906
		//result.Message = err.Error()
		m.writeResult(&oss.Upload{Code: 9906, Message: err.Error()}, writer)
		return
	}

	var file multipart.File
	var fileHeader *multipart.FileHeader
	file, fileHeader, err = request.FormFile("file")
	if err != nil {
		//result.Code = 9905
		//result.Message = err.Error()
		m.writeResult(&oss.Upload{Code: 9905, Message: err.Error()}, writer)
		return
	}

	saveDir := strings.Trim(request.FormValue("path"), "/")
	//override := strings.EqualFold(request.FormValue("override"), "true")
	name := request.FormValue("name")
	if len(name) == 0 {
		name = fileHeader.Filename
		/*fileType := request.FormValue("fileType")
		fileByte, err := io.ReadAll(file)
		if err != nil {
			//result.Code = 9904
			//result.Message = err.Error()
			m.writeResult(&oss.Upload{Code: 9904, Message: err.Error()}, writer)
			return
		}
		name = strings.ToLower(encryption.Md5ByBytes(fileByte)) + fileType*/
	}
	result := m.upload(file, fileHeader, saveDir, name)
	m.writeResult(result, writer)
	log.Println("上传文件", name, fmt.Sprintf("%+v", result))
}

type CheckResult struct {
	FileRootDir, Name string
	SHA256            string
}

func checkAndWrite(newBody []byte, fileRootDir, name string) (*CheckResult, error) {
	ext := filepath.Ext(name)
	path := filepath.Join(fileRootDir, name)

	sha := sha256.New()
	sha.Write(newBody)
	newSha256 := hex.EncodeToString(sha.Sum(nil))

	info, err := os.Stat(path)
	if err == nil && info.Size() > 0 {
		//存在
		readFile, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		sha.Reset()
		sha.Write(readFile)
		oldSHA256 := hex.EncodeToString(sha.Sum(nil))
		if oldSHA256 != newSha256 {
			name = strings.ReplaceAll(name, ext, "-"+time.Now().Format("2006-01-02-150405.999999999")+ext)
			path = filepath.Join(fileRootDir, name)

			err = os.WriteFile(path, newBody, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err = os.WriteFile(path, newBody, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return &CheckResult{
		FileRootDir: fileRootDir,
		Name:        name,
		SHA256:      newSha256,
	}, nil
}
func (m *UploadHandle) upload(file multipart.File, fileHeader *multipart.FileHeader, saveDir string, name string) *oss.Upload {
	fileRootDir := filepath.Join("assets", saveDir) //fmt.Sprintf("assets/%s", path)
	{
		assetsPath, err := filepath.Abs("assets")
		if err != nil {
			//result.Code = 9904
			//result.Message = err.Error()
			return &oss.Upload{Code: 9904, Message: err.Error()}
		}
		savePath, err := filepath.Abs(filepath.Join(fileRootDir, name))
		if len(savePath) < len(assetsPath) || !strings.EqualFold(savePath[:len(assetsPath)], assetsPath) {
			//result.Code = 9904
			//result.Message = "路径不对"
			return &oss.Upload{Code: 9904, Message: "路径不对"}
		}
	}
	if err := os.MkdirAll(fileRootDir, os.ModePerm); err != nil {
		var pathErr *fs.PathError
		ok := errors.As(err, &pathErr)
		if !ok && pathErr != nil {
			//result.Code = 9901
			//result.Message = pathErr.Error()
			return &oss.Upload{Code: 9901, Message: pathErr.Error()}
		}
	}

	openFile, err := fileHeader.Open()
	if err != nil {
		return &oss.Upload{Code: 9901, Message: err.Error()}
	}
	body, err := io.ReadAll(openFile)
	if err != nil {
		return &oss.Upload{Code: 9901, Message: err.Error()}
	}

	var upload = &oss.Upload{Code: 0, Message: "OK"}

	sha := sha256.New()
	sha.Write(body)
	SHA256 := hex.EncodeToString(sha.Sum(nil))

	ext := filepath.Ext(name)

	var result *CheckResult
	//now := time.Now().UnixMilli()
	img, format, err := image.Decode(bytes.NewReader(body))
	//log.Println("image.Decode", time.Now().UnixMilli()-now)
	if err != nil {
		//不是图片或,不支持的图片
		/*{
			//判断是否有一样的数据
			oldFileName := filepath.Join(fileRootDir, name)
			if fileInfo, _ := os.Stat(oldFileName); fileInfo != nil {
				readFile, err := os.ReadFile(oldFileName)
				if err != nil {
					return &oss.Upload{Code: 9901, Message: err.Error()}
				}

				sha.Reset()
				sha.Write(readFile)
				oldSHA256 := hex.EncodeToString(sha.Sum(nil))
				if oldSHA256 == SHA256 {
					//相同直接返回
					upload.Data.Ext = ext
					upload.Data.Format = ""
					upload.Data.SHA256 = SHA256
					upload.Data.Filename = name
					upload.Data.Size = len(body)
					upload.Data.Width = 0
					upload.Data.Height = 0

					urlPath, err := url.JoinPath(saveDir, name)
					if err != nil {
						return &oss.Upload{Code: 9901, Message: err.Error()}
					}
					upload.Data.Path = fmt.Sprintf("/%s", urlPath)
					return upload
				}

				if len(ext) == 0 {
					name = name + "-" + time.Now().Format("20060102150405")
				} else {
					name = strings.ReplaceAll(name, ext, "-"+time.Now().Format("20060102150405")+ext)
				}

			}
		}*/

		result, err = checkAndWrite(body, fileRootDir, name)
		if err != nil {
			return &oss.Upload{Code: 9901, Message: err.Error()}
		}

		upload.Data.Ext = ext
		upload.Data.Format = ""
		upload.Data.SHA256 = SHA256
		upload.Data.Filename = result.Name
		upload.Data.Size = len(body)
		upload.Data.Width = 0
		upload.Data.Height = 0

	} else {
		switch strings.ToUpper(format) {
		case "WEBP":
		case "JPEG":
		case "PNG":
		case "GIF":
		default:
			return &oss.Upload{Code: 9903, Message: fmt.Errorf("不支持图片格式:%s", format).Error()}
		}

		if len(ext) == 0 {
			name = name + "." + format
		}
		ext = filepath.Ext(name)
		name = strings.ReplaceAll(name, ext, ".webp")

		//now = time.Now().UnixMilli()
		//var imgBuf = bytes.NewBuffer(nil)
		//_ = png.Encode(imgBuf, img)
		//log.Println("png.Encode", time.Now().UnixMilli()-now)
		//imgBytes := imgBuf.Bytes()

		var tempFile *os.File
		tempFile, err = os.CreateTemp(os.TempDir(), "oss.*.image")
		if err != nil {
			return &oss.Upload{Code: 9903, Message: err.Error()}
		}
		_, err = tempFile.Write(body)
		if err != nil {
			return &oss.Upload{Code: 9903, Message: err.Error()}
		}
		tempFile.Close()

		var preEncodeFile *os.File
		preEncodeFile, err = os.CreateTemp(os.TempDir(), "oss.*.image")
		if err != nil {
			return &oss.Upload{Code: 9903, Message: err.Error()}
		}
		//改用webp图片格式,如果webp不支持的格式,直接保存
		if strings.ToUpper(format) == "GIF" {
			if err := webpicture.EncodeGIF(tempFile.Name(), preEncodeFile.Name()); err != nil {
				return &oss.Upload{Code: 9903, Message: err.Error()}
			}
		} else {
			if err := webpicture.Encode(tempFile.Name(), preEncodeFile.Name()); err != nil {
				return &oss.Upload{Code: 9903, Message: err.Error()}
			}
		}

		body, err = io.ReadAll(preEncodeFile)
		if err != nil {
			return &oss.Upload{Code: 9903, Message: err.Error()}
		}

		result, err = checkAndWrite(body, fileRootDir, name)
		if err != nil {
			return &oss.Upload{Code: 9901, Message: err.Error()}
		}

		preEncodeFile.Close()

		upload.Data.Ext = ext
		upload.Data.Format = format
		upload.Data.SHA256 = SHA256
		upload.Data.Filename = result.Name
		upload.Data.Size = len(body)
		upload.Data.Width = img.Bounds().Dx()
		upload.Data.Height = img.Bounds().Dy()
	}

	urlPath, err := url.JoinPath(saveDir, result.Name)
	if err != nil {
		return &oss.Upload{Code: 9901, Message: err.Error()}
	}
	upload.Data.Path = fmt.Sprintf("/%s", urlPath)
	return upload
}
func (m *UploadHandle) writeResult(result *oss.Upload, writer http.ResponseWriter) {
	var body []byte
	var err error
	body, err = json.Marshal(result)
	if err != nil {
		log.Println(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = writer.Write(body)
	if err != nil {
		log.Println(err)
	}
}
