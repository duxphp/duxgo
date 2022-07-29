package form

import (
	"github.com/duxphp/duxgo/ui/node"
)

type FileType string

const (
	FileUpload FileType = "upload"
	FileManage FileType = "manage"
)

// File 文件上传
type File struct {
	mode      FileType
	url       string
	manageUrl string
}

// NewFile 创建文本
func NewFile() *File {
	return &File{}

}

// Type 上传模式
func (a *File) Type(mode FileType) *File {
	a.mode = mode
	return a
}

// Url 上传地址
func (a *File) Url(url string) *File {
	a.url = url
	return a
}

// ManageUrl 管理地址
func (a *File) ManageUrl(manageUrl string) *File {
	a.manageUrl = manageUrl
	return a
}

// GetValue 格式化值
func (a *File) GetValue(value any, info map[string]any) any {
	return value
}

// SaveValue 保存数据
func (a *File) SaveValue(value any, data map[string]any) any {
	return value
}

// Render 渲染
func (a *File) Render(element node.IField) *node.TNode {
	ui := node.TNode{
		"nodeName":     "app-file",
		"upload":       a.url,
		"fileUrl":      a.manageUrl,
		"type":         a.mode,
		"vModel:value": element.GetUIField(),
	}
	return &ui
}
