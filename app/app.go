package app

var (
	AppList  = make(map[string]*AppConfig)
	AppIndex []string
)

// AppConfig 注册规则
type AppConfig struct {
	Name     string //应用名称
	Config   any    //应用配置
	Init     func() //初始化应用
	Register func() // 注册应用
	Boot     func() // 启动应用
}

// App 注册应用
func App(opt *AppConfig) {
	AppList[opt.Name] = opt
	AppIndex = append(AppIndex, opt.Name)
}
