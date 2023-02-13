package views

import (
	"encoding/json"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

func Init() {
	// 注册模板引擎
	funcMap := template.FuncMap{
		"unescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"marshal": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	}
	tpl := template.Must(template.New("").Delims("${", "}").Funcs(funcMap).ParseFS(registry.TplFs, "template/*"))
	registry.Tpl = tpl
}

// Template 模板服务
type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
