package service

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"

	"dandelion/app/service/dao"
	"dandelion/app/util"
)

type FileService struct{}

func (self FileService) LoadAction(context *gweb.Context) gweb.Result {
	path := context.Request.URL.Query().Get("path")

	return &gweb.ImageResult{FilePath: path}
}

func (self FileService) WriteFile(b []byte, ContentType string) string {
	md5Name := util.Md5ByBytes(b)

	now := time.Now()

	var f *os.File

	fileType := strings.Split(ContentType, "/")[1]
	fileType = strings.Split(fileType, "+")[0]
	filePath := "upload/" + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day()) + "/" + md5Name[0:2]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, os.ModePerm)
		tool.CheckError(err)
	}

	fileName := filePath + "/" + md5Name + "." + fileType

	if _, err := os.Stat(fileName); os.IsNotExist(err) {

		f, err = os.Create(fileName) //创建文件
		tool.CheckError(err)
		defer f.Close()
		f.Write(b)
		f.Sync()

	} else {
		//f, err = os.OpenFile(fileName, os.O_RDONLY, os.ModePerm) //打开文件
		//tool.CheckError(err)
		//fmt.Println(f)
	}
	return fileName
}
func (self FileService) UploadAction(context *gweb.Context) gweb.Result {
	context.Request.ParseForm()
	File, FileHeader, err := context.Request.FormFile("file")
	tool.CheckError(err)
	b, err := ioutil.ReadAll(File)
	tool.CheckError(err)
	defer File.Close()

	fileName := self.WriteFile(b, FileHeader.Header.Get("Content-Type"))
	//base64Data := "data:" + FileHeader.Header.Get("Content-Type") + ";base64," + base64.StdEncoding.EncodeToString(b)

	//framework.WriteJSON(context, &framework.ActionStatus{true, "oK", base64Data})
	return &gweb.JsonResult{Data: &dao.ActionStatus{true, "ok", fileName}}
}

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

	return self.WriteFile(b, resp.Header.Get("Content-Type"))

}
