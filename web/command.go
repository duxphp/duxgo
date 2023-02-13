package web

import (
	"github.com/duxphp/duxgo/v2/app"
	"github.com/duxphp/duxgo/v2/service"
	"github.com/duxphp/duxgo/v2/task"
	"github.com/spf13/cobra"
)

func Command(command *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "starting the web service",
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化服务
			service.Init()
			task.Init()
			Init()
			// 初始化应用
			app.Init()
			// 启动web服务
			Start()
		},
	}
	command.AddCommand(cmd)
}
