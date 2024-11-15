package httpext

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/funcmap"
)

type DefaultViewRender struct {
	ViewDir string
}

func (v *DefaultViewRender) Render(context constrain.IContext, request *http.Request, writer http.ResponseWriter, viewData constrain.IViewResult) error {
	if len(v.ViewDir) == 0 {
		v.ViewDir = "view"
	}
	var err error
	var fileByte []byte

	//vd := viewData.(constrain.IViewResult)

	var errs []error

	viewName := viewData.GetName()
	if len(viewName) > 0 {
		dir, _ := filepath.Split(context.Route())
		if dir == "/" {
			fileByte, err = os.ReadFile(fmt.Sprintf("%s/%s.html", v.ViewDir, viewName))
			if err != nil {
				if err!=nil{
					errs =append(errs,err)
				}
				fileByte = []byte(err.Error())
				err = nil
			}
		} else {
			fileByte, err = os.ReadFile(fmt.Sprintf("%s%s%s.html", v.ViewDir, dir, viewName))
			if err!=nil{
				errs =append(errs,err)
			}
			if err != nil {
				fileByte, err = os.ReadFile(fmt.Sprintf("%s/404.html", v.ViewDir))
				if err!=nil{
					errs =append(errs,err)
				}
			}
		}
	} else {
		path := strings.Trim(context.Route(), "/")
		ext := filepath.Ext(path)
		if len(ext) > 0 {
			fileByte, err = os.ReadFile(fmt.Sprintf("%s/%s", v.ViewDir, path))
			if err!=nil{
				errs =append(errs,err)
			}
		} else {
			fileByte, err = os.ReadFile(fmt.Sprintf("%s/%s.html", v.ViewDir, path))
			if err != nil {
				errs =append(errs,err)
				fileByte, err = os.ReadFile(fmt.Sprintf("%s/%s.html", v.ViewDir, "index"))
				if err!=nil{
					errs =append(errs,err)
				}
			}
		}
	}

	if err != nil {

		return errors.Join(errs...)
	}

	var t *template.Template
	t, err = template.New("").Funcs(funcmap.NewFuncMap().Build(context)).Parse(string(fileByte))
	if err != nil {
		return err
	}

	filenames, err := filepath.Glob(fmt.Sprintf("%s/template/*.gohtml", v.ViewDir))
	if err != nil {
		return err
	}
	if len(filenames) > 0 {
		t, err = t.ParseFiles(filenames...)
		if err != nil {
			return err
		}
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	contextValue := contexext.FromContext(context)
	err = t.Execute(writer, map[string]interface{}{
		"Query": contextValue.Query,
		"Data":  viewData,
	})
	if err != nil {
		//return err
		writer.WriteHeader(http.StatusNotFound)
		_, err = writer.Write([]byte(err.Error()))
		if err != nil {
			return err
		}
	}
	return nil
}
