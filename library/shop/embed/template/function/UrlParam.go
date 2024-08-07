package function

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/funcmap"
	"github.com/nbvghost/dandelion/library/util"
	"strings"
)

type UrlParam struct {
	Type   string   `arg:""` //add del clear
	Values []string `arg:"..."`
}

func (g *UrlParam) Call(ctx constrain.IContext) funcmap.IFuncResult {
	if len(g.Values)%2 != 0 {
		return funcmap.NewStringFuncResult("参数必须是偶数个")
	}
	var contextValue = contexext.FromContext(ctx)
	var u = util.GetHost(contextValue.Request) + contextValue.Request.URL.Path

	params := contextValue.Request.URL.Query()
	var label string
	for i := range g.Values {
		if (i+1)%2 == 1 {
			label = strings.ToLower(g.Values[i])
		} else {

			switch g.Type {
			case "set":
				var newValue = strings.ToLower(g.Values[i])
				params.Set(label, newValue)
			case "add":
				{
					var hasValue = params.Get(label)
					if hasValue == "" {
						params.Add(label, strings.ToLower(g.Values[i]))
					} else {
						var newValue = strings.ToLower(g.Values[i])
						if strings.EqualFold(hasValue, newValue) {
							//params.set(label, newValue)
						} else {
							params.Add(label, newValue)
						}
					}
				}
			case "del":
				params.Del(label)
			}
		}
	}
	switch g.Type {
	case "clear":
		if len(g.Values) > 0 {
			for key := range params {
				params.Del(key)
			}
		}
	}
	return funcmap.NewStringFuncResult(u + "?" + params.Encode())
}
