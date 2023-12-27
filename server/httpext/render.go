package httpext

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
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

	viewName := viewData.GetName()
	if len(viewName) > 0 {
		dir, _ := filepath.Split(context.Route())
		if dir == "/" {
			fileByte, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.html", v.ViewDir, viewName))
			if err != nil {
				fileByte = []byte(err.Error())
				err = nil
			}
		} else {
			fileByte, err = ioutil.ReadFile(fmt.Sprintf("%s%s%s.html", v.ViewDir, dir, viewName))
			if err != nil {
				fileByte, err = ioutil.ReadFile(fmt.Sprintf("%s/404.html", v.ViewDir))
			}
		}
	} else {
		path := strings.Trim(context.Route(), "/")
		ext := filepath.Ext(path)
		if len(ext) > 0 {
			fileByte, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", v.ViewDir, path))
		} else {
			fileByte, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.html", v.ViewDir, path))
			if err != nil {
				if _, ok := err.(*fs.PathError); ok {
					fileByte, err = ioutil.ReadFile(fmt.Sprintf("%s/%s.html", v.ViewDir, "index"))
				}
			}
		}
	}

	if err != nil {
		return err
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
		return err
	}
	return nil
}
