package httpext

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/funcmap"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type viewRender struct {
}

func (v *viewRender) Render(context constrain.IContext, request *http.Request, writer http.ResponseWriter, viewData interface{}) error {
	var err error
	var fileByte []byte

	vd := viewData.(constrain.IViewResult)

	viewName := vd.GetName()
	if len(viewName) > 0 {
		dir, _ := filepath.Split(context.Route())
		if dir == "/" {
			fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s.html", viewName))
			if err != nil {
				fileByte = []byte(err.Error())
				err = nil
			}
		} else {
			fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s/%s.html", dir, viewName))
		}

	} else {
		fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s.html", strings.TrimSuffix(context.Route(), "/")))
		if err != nil {
			if _, ok := err.(*fs.PathError); ok {
				fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s.html", "index"))
			}
		}

	}

	if err != nil {
		return err
	}

	var t *template.Template
	t, err = template.New("").Funcs(funcmap.NewFuncMap(context)).Parse(string(fileByte))
	if err != nil {
		return err
	}

	filenames, err := filepath.Glob(fmt.Sprintf("view/template/*.gohtml"))
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
