package web

import (
	"github.com/duxphp/duxgo/v2/app"
	"github.com/duxphp/duxgo/v2/global"
	"github.com/duxphp/duxgo/v2/logger"
	"github.com/duxphp/duxgo/v2/monitor"
	"github.com/duxphp/duxgo/v2/service"
	"github.com/duxphp/duxgo/v2/task"
	"github.com/duxphp/duxgo/v2/websocket"
	"github.com/gookit/color"
	"github.com/gookit/event"
	"github.com/panjf2000/ants/v2"
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
			signal.Notify(ch,
				os.Interrupt,
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM)

			service.Init()
			task.Init()
			Init()
			monitor.Init()
			app.Init()

			task.ListenerTask("dux.monitor", monitor.Control)
			task.ListenerScheduler("*/1 * * * *", "dux.monitor", map[string]any{}, task.PRIORITY_LOW)
			// Start timing service
			go func() {
				task.StartScheduler()
			}()
			// Start queue service
			go func() {
				task.Add("ping", &map[string]any{})
				task.StartQueue()
			}()
			// Starting the web service
			go func() {
				Start()
			}()
			<-ch
			// Shut down service
			color.Println("⇨ <orange>Stop scheduler</>")
			task.StopScheduler()
			color.Println("⇨ <orange>Stop queue</>")
			task.StopQueue()
			color.Println("⇨ <orange>Stop event</>")
			err, _ := event.Fire("app.close", event.M{})
			if err != nil {
				logger.Log().Error().Err(err).Msg("event stop")
			}
			color.Println("⇨ <orange>Stop websocket</>")
			websocket.Release()
			color.Println("⇨ <orange>Stop pools</>")
			ants.Release()
			color.Println("⇨ <orange>Stop fiber</>")
			_ = global.App.ShutdownWithTimeout(0)
			color.Println("⇨ <red>Server closed</>")
		},
	}
	command.AddCommand(cmd)
}
