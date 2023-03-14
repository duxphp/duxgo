<p align="center">
<a href="https://www.duxravel.com/">
    <img src="https://github.com/duxphp/duxravel/blob/main/resources/image/watermark.png?raw=true" width="100" height="100">
</a>

<p align="center"><code>duxgo</code> 是一款基于 GoFiber 的快速开发框架，集成主流三方包，简单、易开发、高性能的集成框架。</p>

<p align="center">
<a href="https://www.duxfast.com">中文文档</a>
</p>


# 🎯 特点

- 📦 基于 GoFiber 的 Fasthttp 高性能 Web 框架。
- 📚 整合 Gorm 作为主要数据库驱动，提供良好的数据库操作支持。
- 📡 不做过度封装，便于开发者灵活选择和随版本升级。
- 🔧 集成各大流行包，并封装常用日志、异常、权限等工具包。
- 📡 采用应用模块化设计，提高应用程序的可维护性和可扩展性。
- 📡 统一注册应用入口，方便应用程序的整体架构和管理。
- 🏷 开发命令助手与脚手架工具，提供基础的代码生成。


#  ⚡ 快速开始

```go
package main

import (
	"github.com/duxphp/duxgo/v2/app"
	"github.com/duxphp/duxgo/v2/route"
)

func main() {
	dux := duxgo.New()
	
	app := route.Add("web", route.New(""))

	app.Get("/", func(c *fiber.Ctx) error {
		return  c.SendString("Hello, World 👋!")
	}, "首页", "web.home")
	
	dux.Run()
}

```

#  ⚙ 安装

请确保当前 Golang 环境版本高于 `1.18` 版本，建立项目目录并初始化。

```sh
go get github.com/duxphp/duxgo/v2
```

# 💡思想

该框架遵循与 DuxLite 一致化架构设计，将各个功能模块应用化，并通过 `应用入口` 与 `事件调度` 进行高度解耦，并保证基础框架与系统必备最小化，避免大而全的臃肿框架设计。