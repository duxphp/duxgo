package web

import (
	"github.com/duxphp/duxgo/v2/app"
	"github.com/duxphp/duxgo/v2/monitor"
	"github.com/duxphp/duxgo/v2/service"
	"github.com/duxphp/duxgo/v2/task"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func Command(command *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "starting the web service",
		Run: func(cmd *cobra.Command, args []string) {

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt,
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM)

			// 初始化服务
			service.Init()
			// 初始化任务
			task.Init()
			// 初始化WEB
			Init()
			// 初始化监控
			monitor.Init()
			// 初始化应用
			app.Init()
			// 注册监控
			task.RegTask("dux.monitor", monitor.Control)
			task.RegScheduler("*/1 * * * *", "dux.monitor", map[string]any{}, task.PRIORITY_LOW)
			// 启动定时服务
			go func() {
				task.StartScheduler()
			}()
			// 启动队列服务
			go func() {
				task.Add("ping", &map[string]any{})
				task.StartQueue()
			}()
			// 启动web服务
			Start()
			<-ch
			// 关闭服务
			task.StopScheduler()
			task.StopQueue()
			Stop()
			color.Println("⇨ <red>Server closed</>")
		},
	}
	command.AddCommand(cmd)
}
