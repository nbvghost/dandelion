package result

type ViewResult interface {
	GetName() string
	SetName(name string)
}

/*func (r *HtmlResult) Render(context context.IContext) ([]byte, error) {
	var err error
	var fileByte []byte
	fileByte, err = ioutil.ReadFile(fmt.Sprintf("view/%s.html", strings.TrimSuffix(context.Route(), "/")))
	if err != nil {
		return nil, err
	}
	var t *template.Template
	t, err = template.New("").Funcs(funcmap.NewFuncMap(context)).Parse(string(fileByte))
	if err != nil {
		return nil, err
	}
	t, err = t.ParseGlob("view/template/*.gohtml")
	if err != nil {
		return nil, err
	}

	r.Data["Query"] = context.Query()

	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, r.Data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil

}
*/
