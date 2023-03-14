package views

import (
	"embed"
	"encoding/json"
	"github.com/gofiber/template/html"
	"html/template"
	"net/http"
)

var TplFs embed.FS

var FrameFs embed.FS

var FrameTpl *template.Template

var Views *html.Engine

func Init() {
	// 注册框架模板
	FrameTpl = template.Must(template.New("").ParseFS(FrameFs, "template/*"))

	// 注册 Fiber 引擎
	engine := html.NewFileSystem(http.FS(TplFs), ".gohtml")
	engine.AddFunc("unescape", func(v string) template.HTML {
		return template.HTML(v)
	})
	engine.AddFunc("marshal", func(v string) string {
		a, _ := json.Marshal(v)
		return string(a)
	})
	Views = engine
}
