duxgo
====================

# 概述
duxgo 是一款基于 go-echo 框架整合常用的 ORM、日志、队列、缓存等 web 开发常用功能，提供了一个简单、易用、灵活的框架。

# 安装

使用 go get 安装，无任何第三方依赖：

```sh
go get github.com/duxphp/duxgo
```

# 使用方法

## 1. 创建服务

```go
import "github.com/duxphp/duxgo"
server := duxgo.New()
```

## 2. 设置配置目录

将含有 toml 配置文件的目录添加到配置目录中
```go
server.SetConfigDir("./config")
```

## 3. 注册应用

将自己开发的 HMVC 结构的应用注册到框架中
```go
server.Register(func(app *duxgo.App) {
    system.App()
    tools.App()
})
```

## 4. 启动服务
```go
server.Start()
```
