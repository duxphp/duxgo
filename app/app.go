package app

var (
	// List 应用列表
	List = make(map[string]*Config)

	// Indexes 应用索引
	Indexes []string
)

// Config 注册规则
type Config struct {
	Name     string //应用名称
	Config   any    //应用配置
	Title    string //应用标题
	Desc     string //应用描述
	Init     func() //初始化应用
	Register func() // 注册应用
	Boot     func() // 启动应用
}

// Register 注册应用
func Register(opt *Config) {
	List[opt.Name] = opt
	Indexes = append(Indexes, opt.Name)
}
