package file

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
)

type NoneUp struct {
}

func (m *NoneUp) HandleOptions(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *NoneUp) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//m.fileUploadAction(context, "none")

	return &result.EmptyResult{}, err
}

func (m *NoneUp) HandleGet(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *NoneUp) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

/*
func (m *NoneUp) fileUploadAction(context constrain.IContext, dynamicDirName string) {

	context.Request.ParseForm()
	File, FileHeader, err := context.Request.FormFile("file")
	if err != nil {
		result := make(map[string]interface{})
		result["Success"] = false
		result["Message"] = err
		result["Path"] = ""
		result["Url"] = ""
		rb, _ := json.Marshal(result)
		context.Response.Write(rb)
		return
	}
	defer File.Close()

	err, fileName := gweb.WriteWithFile(File, FileHeader, dynamicDirName, "images")
	if err != nil {
		result := make(map[string]interface{})
		result["Success"] = false
		result["Message"] = err
		result["Path"] = ""
		//result["Url"] = ""
		rb, _ := json.Marshal(result)
		context.Response.Write(rb)
	} else {
		result := make(map[string]interface{})
		result["Success"] = true
		result["Message"] = "OK"
		result["Path"] = fileName
		//result["Url"] = "//" + conf.Config.Domain + "/file/load?path=" + fileName
		rb, _ := json.Marshal(result)
		context.Response.Write(rb)
	}

}*/
