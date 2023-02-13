package task

import (
	"github.com/duxphp/duxgo/v2/service"
	"github.com/spf13/cobra"
)

func Command(command *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "Start the queue and schedule tasks",
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化服务
			service.Init()
			Init()
			//启动队列与调度服务
			go func() {
				StartScheduler()
			}()
			Add("ping", &map[string]any{})
			StartQueue()
		},
	}
	command.AddCommand(cmd)
}
