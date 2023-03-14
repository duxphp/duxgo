<p align="center">
<a href="https://www.duxravel.com/">
    <img src="https://github.com/duxphp/duxravel/blob/main/resources/image/watermark.png?raw=true" width="100" height="100">
</a>
<p align="center"><code>duxgo</code> is a fast development framework based on GoFiber, integrated with mainstream third-party packages, simple, easy to develop, and high-performance integrated framework.</p>
<p align="center">
<a href="https://www.duxfast.com">Development documen</a>
</p>


# 🎯 Features

- 📦 High-performance Web framework based on GoFiber Fasthttp.
- 📚 Integrated Gorm as the main database driver to provide good database operation support.
- 📡 Not overly encapsulated, making it easy for developers to flexibly choose and upgrade with the version.
- 🔧 Integrating major popular packages and encapsulating commonly used tool packages such as logs, exceptions, and permissions.
- 📡 Adopt an application modular design to improve the maintainability and scalability of the application.
- 📡 Uniform registration of application entry points, facilitating the overall architecture and management of the application.
- 🏷 Developing command assistants and scaffolding tools, providing basic code generation.


#  ⚡ Quick start

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
	}, "index", "web.home")
	
	dux.Run()
}

```

#  ⚙ Installation

Please make sure that the current Golang environment version is higher than `1.18`, create the project directory and initialize it.

```sh
go get github.com/duxphp/duxgo/v2
```

# 💡Philosophy

This framework follows the architectural design of DuxLite, applying each functional module and highly decoupling through `application entry points` and `event scheduling`, ensuring the minimization of basic framework and system required to avoid cumbersome framework designs.