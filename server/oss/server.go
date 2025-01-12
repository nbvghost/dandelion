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
	"golang.org/x/image/draw"
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
	"os"
	"path/filepath"
	"strconv"
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

	dir, fileName := filepath.Split(request.URL.Path)
	if strings.Contains(fileName, "@") {
		w := 0
		h := 0
		fileNames := strings.Split(fileName, "@")
		fileName = fileNames[0]
		//gIndex := strings.LastIndex(fileNames[0], "@")
		if len(fileNames) >= 2 {
			ss := fileNames[1]
			wh := strings.Split(ss, "x")
			if len(wh) == 1 {
				w, _ = strconv.Atoi(wh[0])
				h, _ = strconv.Atoi(wh[0])
			} else if len(wh) == 2 {
				w, _ = strconv.Atoi(wh[0])
				h, _ = strconv.Atoi(wh[1])
				if w == 0 || h == 0 {
					w = w + h
					h = w
				}
			}
		}

		if w > 0 && h > 0 {
			{
				//读取临时文件
				if f, err := os.Open(filepath.Join(os.TempDir(), "oss", request.URL.Path)); err != nil {
					log.Println(err)
				} else {
					defer func(f *os.File) {
						err := f.Close()
						if err != nil {
							log.Println(err)
						}
					}(f)
					fullNameFileInfo, err := f.Stat()
					if err != nil {
						writeToErrImage(writer, request, err)
						return
					} else {
						if fullNameFileInfo.IsDir() {
							writeToDefaultImage(writer)
							return
						} else {
							b, err := io.ReadAll(f)
							if err != nil {
								writeToErrImage(writer, request, err)
								return
							}

							imgFile := bytes.NewReader(b)

							img, s, err := image.DecodeConfig(imgFile)
							if err != nil {
								writeToErrImage(writer, request, err)
								return
							}
							writeToImage(writer, b, "image/"+s, img.Width, img.Height)
							return
						}
					}
				}
			}
			var imgFile io.Reader
			if f, err := os.Open(filepath.Join("assets", dir, fileName)); err != nil {
				writeToErrImage(writer, request, err)
				return
			} else {
				fullNameFileInfo, err := f.Stat()
				if err != nil {
					writeToErrImage(writer, request, err)
					return
				} else {
					if fullNameFileInfo.IsDir() {
						writeToDefaultImage(writer)
						return
					} else {
						imgFile = f
					}
				}
			}
			img, s, err := image.Decode(imgFile)
			if err != nil {
				writeToErrImage(writer, request, err)
				return
			}
			//p := float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y)
			if float64(w)/float64(img.Bounds().Max.X) > float64(h)/float64(img.Bounds().Max.Y) {
				//w = int(float64(h) * p)
				w = int(float64(img.Bounds().Max.X) * (float64(h) / float64(img.Bounds().Max.Y)))
			} else {
				h = int(float64(img.Bounds().Max.Y) * (float64(w) / float64(img.Bounds().Max.X)))
			}

			dst := image.NewRGBA(image.Rect(0, 0, w, h))
			draw.ApproxBiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Src, nil)

			buffer := bytes.NewBuffer(nil)
			switch s {
			case "jpeg":
				err = jpeg.Encode(buffer, dst, nil)
			case "png":
				err = png.Encode(buffer, dst)
			case "gif":
				err = gif.Encode(buffer, dst, &gif.Options{})
			case "bmp":
				err = bmp.Encode(buffer, dst)
			default:
				//err = fmt.Errorf("ERROR FORMAT:%s", s)
				err = png.Encode(buffer, dst)
			}
			if err != nil {
				writeToErrImage(writer, request, err)
				return
			}
			/*writer.WriteHeader(http.StatusOK)
			writer.Header().Set("content-type", "image/"+s)
			writer.Header().Set("width", fmt.Sprintf("%d", dst.Rect.Dx()))
			writer.Header().Set("height", fmt.Sprintf("%d", dst.Rect.Dy()))
			_, err = writer.Write(buffer.Bytes())
			if err != nil {
				log.Println(err)
				return
			}*/

			//写入临时目录，/var/tmp,不用每次都去调整尺寸filepath.Join(dir, fileName)
			tmpDir := filepath.Join(os.TempDir(), "oss", dir)
			_ = os.MkdirAll(tmpDir, os.ModePerm)
			err = os.WriteFile(filepath.Join(os.TempDir(), "oss", request.URL.Path), buffer.Bytes(), os.ModePerm)
			if err != nil {
				writeToErrImage(writer, request, err)
				return
			}

			writeToImage(writer, buffer.Bytes(), "image/"+s, dst.Rect.Dx(), dst.Rect.Dy())
		} else {
			var imgFile io.Reader

			if f, err := os.Open(filepath.Join("assets", dir, fileName)); err != nil {
				writeToErrImage(writer, request, err)
				return
			} else {
				fullNameFileInfo, err := f.Stat()
				if err != nil {
					writeToErrImage(writer, request, err)
					return
				} else {
					if fullNameFileInfo.IsDir() {
						writeToDefaultImage(writer)
						return
					} else {
						imgFile = f
					}
				}
			}

			c, s, err := image.DecodeConfig(imgFile)
			if err != nil {
				http.NotFoundHandler().ServeHTTP(writer, request)
				return
			}
			b, err := io.ReadAll(imgFile)
			if err != nil {
				http.NotFoundHandler().ServeHTTP(writer, request)
				return
			}

			writeToImage(writer, b, "image/"+s, c.Width, c.Height)
		}

	} else {
		fp := filepath.Join("assets", request.URL.Path)
		//log.Println("fp", fp)
		fileInfo, err := os.Stat(fp)
		if os.IsNotExist(err) || fileInfo.IsDir() {
			writeToDefaultImage(writer)
		} else {
			http.FileServer(http.Dir("assets")).ServeHTTP(writer, request)
		}
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
	return &Server{etcdService: etcdService, engine: engine}
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
	override := strings.EqualFold(request.FormValue("override"), "true")
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
	result := m.upload(file, fileHeader, saveDir, name, override)
	m.writeResult(result, writer)
	log.Println("上传文件", name, result.Data.Path)
}

func (m *UploadHandle) upload(file multipart.File, fileHeader *multipart.FileHeader, saveDir string, name string, override bool) *oss.Upload {
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
		pathErr, ok := err.(*fs.PathError)
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
	{
		oldFileName := filepath.Join(fileRootDir, name)
		if fileInfo, _ := os.Stat(oldFileName); fileInfo != nil {
			readFile, err := os.ReadFile(oldFileName)
			if err != nil {
				return nil
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
				upload.Data.Path = fmt.Sprintf("/%s", filepath.Join(saveDir, name))
				return upload
			}

			if len(ext) == 0 {
				name = name + "-" + time.Now().Format("20060102150405")
			} else {
				name = strings.ReplaceAll(name, ext, "-"+time.Now().Format("20060102150405")+ext)
			}

		}
	}

	now := time.Now().UnixMilli()
	img, format, err := image.Decode(bytes.NewReader(body))
	log.Println("image.Decode", time.Now().UnixMilli()-now)
	if err != nil {
		err = os.WriteFile(filepath.Join(fileRootDir, name), body, os.ModePerm)
		if err != nil {
			return &oss.Upload{Code: 9903, Message: err.Error()}
		}
		upload.Data.Ext = ext
		upload.Data.Format = ""
		upload.Data.SHA256 = SHA256
		upload.Data.Filename = name
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
		now = time.Now().UnixMilli()
		//imgBytes := imgBuf.Bytes()
		err = os.WriteFile(filepath.Join(fileRootDir, name), body, os.ModePerm)
		if err != nil {
			return &oss.Upload{Code: 9903, Message: err.Error()}
		}
		log.Println("os.WriteFile", time.Now().UnixMilli()-now)

		now = time.Now().UnixMilli()
		//改用webp图片格式,如果webp不支持的格式,直接保存
		if strings.ToUpper(format) == "GIF" {
			if err := webpicture.EncodeGIF(filepath.Join(fileRootDir, name), filepath.Join(fileRootDir, name)); err != nil {
				return &oss.Upload{Code: 9903, Message: err.Error()}
			}
		} else {
			if err := webpicture.Encode(filepath.Join(fileRootDir, name), filepath.Join(fileRootDir, name)); err != nil {
				return &oss.Upload{Code: 9903, Message: err.Error()}
			}
		}

		log.Println("webpicture.Encode", time.Now().UnixMilli()-now)

		upload.Data.Ext = ext
		upload.Data.Format = format
		upload.Data.SHA256 = SHA256
		upload.Data.Filename = name
		upload.Data.Size = len(body)
		upload.Data.Width = img.Bounds().Dx()
		upload.Data.Height = img.Bounds().Dy()
	}

	upload.Data.Path = fmt.Sprintf("/%s", filepath.Join(saveDir, name))
	return upload

	/*if len(saveDir) == 0 {
		upload.Data.Path = fmt.Sprintf("/%s", name)
		upload.Data.Ext = filepath.Ext(filepath.Ext(fileHeader.Filename))
		return upload
	} else {
		upload.Data.Path = fmt.Sprintf("/%s", filepath.Join(saveDir, name))
		upload.Data.Ext = filepath.Ext(filepath.Ext(fileHeader.Filename))
		return upload
	}*/
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
