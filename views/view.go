package views

import (
	"embed"
	"encoding/json"
	"github.com/gofiber/template/html"
	"github.com/samber/do"
	"html/template"
	"net/http"
)

var TplFs embed.FS

func Tpl() *html.Engine {
	return do.MustInvoke[*html.Engine](nil)
}

func Init() {
	// 注册模板引擎
	engine := html.NewFileSystem(http.FS(TplFs), ".gohtml")
	engine.AddFunc("unescape", func(v string) template.HTML {
		return template.HTML(v)
	})
	engine.AddFunc("marshal", func(v string) string {
		a, _ := json.Marshal(v)
		return string(a)
	})
	do.ProvideValue[*html.Engine](nil, engine)
}
