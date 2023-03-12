package views

import (
	"embed"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"html/template"
	"io"
)

var TplFs embed.FS

func Tpl() *template.Template {
	return do.MustInvoke[*template.Template](nil)
}

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
	do.ProvideValue[*template.Template](nil, template.Must(template.New("").Delims("${", "}").Funcs(funcMap).ParseFS(TplFs, "template/*")))
}

// Template 模板服务
type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
