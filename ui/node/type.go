package node

type (
	TNode map[string]any
)

// IField 表单元素接口
type IField interface {
	GetUIField(field ...string) string
	GetName() string
}

type IForm interface {
	AddField(name string, field string)
}

type IWidget interface {
	Render() *TNode
}
