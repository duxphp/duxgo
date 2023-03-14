package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duxphp/duxgo/v2/config"
	"github.com/duxphp/duxgo/v2/global"
	"github.com/duxphp/duxgo/v2/logger"
	"github.com/gookit/color"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"github.com/spf13/cast"
	"time"
)

// Init 初始化任务处理
func Init() {
	dbConfig := config.Get("database").GetStringMapString("redis")
	res := asynq.RedisClientOpt{
		Addr: dbConfig["host"] + ":" + dbConfig["port"],
		// Omit if no password is required
		Password: dbConfig["password"],
		// Use a dedicated db number for asynq.
		// By default, Redis offers 16 databases (0..15)
		DB: cast.ToInt(dbConfig["db"]),
	}

	// 普通队列服务
	srv := asynq.NewServer(
		res,
		asynq.Config{
			Logger: &TaskLogger{
				Logger: logger.New(
					logger.GetWriter(
						zerolog.LevelDebugValue,
						"task",
						"default",
						true,
					),
				).With().Timestamp().Logger(),
			},
			LogLevel:    asynq.WarnLevel,
			Concurrency: 20,
			Queues: map[string]int{
				"high":    10,
				"default": 7,
				"low":     3,
			},
			//ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			//	retried, _ := asynq.GetRetryCount(ctx)
			//	maxRetry, _ := asynq.GetMaxRetry(ctx)
			//	if retried >= maxRetry {
			//		logger.Log().Info().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("task retry")
			//	} else {
			//		logger.Log().Info().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("task error")
			//	}
			//}),
		},
	)

	// 混合器
	mux := asynq.NewServeMux()

	// 客户端
	client := asynq.NewClient(res)

	// 检查器
	inspector := asynq.NewInspector(res)

	// 队列服务端
	do.ProvideValue[*asynq.Server](nil, srv)
	// 队列混合器
	do.ProvideValue[*asynq.ServeMux](nil, mux)
	// 队列客户端
	do.ProvideValue[*asynq.Client](nil, client)
	// 队列检查器
	do.ProvideValue[*asynq.Inspector](nil, inspector)

	mux.HandleFunc("ping", func(ctx context.Context, t *asynq.Task) error {
		color.Print("⇨ <green>Task server start</>\n")
		return nil
	})

	// 定时调度服务
	scheduler := asynq.NewScheduler(res, &asynq.SchedulerOpts{
		LogLevel: asynq.ErrorLevel,
		Location: global.TimeLocation,
		PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
			if err == nil {
				return
			}
			logger.Log().Error().Msgf("scheduler: ", err.Error())
		},
	})

	// 定时调度器
	do.ProvideValue[*asynq.Scheduler](nil, scheduler)
}

type Priority string

const (
	PRIORITY_HIGH    Priority = "high"
	PRIORITY_DEFAULT Priority = "default"
	PRIORITY_LOW     Priority = "low"
)

// StartQueue 启动队列服务
func StartQueue() {
	if err := do.MustInvoke[*asynq.Server](nil).Run(do.MustInvoke[*asynq.ServeMux](nil)); err != nil {
		logger.Log().Error().Msgf("Queue service cannot be started: %v", err)
	}
	do.MustInvoke[*asynq.Server](nil).Shutdown()
}

// StartScheduler 启动调度服务
func StartScheduler() {
	if err := do.MustInvoke[*asynq.Scheduler](nil).Run(); err != nil {
		logger.Log().Error().Msgf("Scheduler service cannot be started: %v", err)
	}
	do.MustInvoke[*asynq.Scheduler](nil).Shutdown()
}

func StopQueue() {
	do.MustInvoke[*asynq.Server](nil).Shutdown()
}

func StopScheduler() {
	do.MustInvoke[*asynq.Scheduler](nil).Shutdown()

}

//

// Add 添加即时队列
func Add(typename string, params any, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return AddTask(typename, params, asynq.Queue(string(group)))
}

// AddDelay 延迟队列（秒）
func AddDelay(typename string, params any, t time.Duration, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return AddTask(typename, params, asynq.ProcessIn(t), asynq.Queue(string(group)))
}

// AddTime 定时队列（秒）
func AddTime(typename string, params any, t time.Time, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return AddTask(typename, params, asynq.ProcessAt(t), asynq.Queue(string(group)))
}

// AddTask 添加队列任务
func AddTask(typename string, params any, opts ...asynq.Option) *asynq.TaskInfo {
	payload, _ := json.Marshal(params)
	task := asynq.NewTask(typename, payload)
	opts = append(opts, asynq.MaxRetry(3))            // 重试3次
	opts = append(opts, asynq.Timeout(1*time.Minute)) // 1分钟超时
	opts = append(opts, asynq.Retention(2*time.Hour)) // 保留2小时

	info, err := do.MustInvoke[*asynq.Client](nil).Enqueue(task, opts...)
	if err != nil {
		logger.Log().Error().Msg("Queue add error :" + err.Error())
	}
	return info
}

// DelTask 删除队列任务
func DelTask(priority Priority, id string) error {
	err := do.MustInvoke[*asynq.Inspector](nil).DeleteTask(string(priority), id)
	if errors.Is(err, asynq.ErrQueueNotFound) {
		return nil
	}
	if errors.Is(err, asynq.ErrTaskNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// RegScheduler 注册定时任务
func RegScheduler(cron string, typename string, params any, priority ...Priority) {
	payload, _ := json.Marshal(params)
	task := asynq.NewTask(typename, payload)
	var opts []asynq.Option
	opts = append(opts, asynq.MaxRetry(3))             // 重试3次
	opts = append(opts, asynq.Timeout(30*time.Minute)) // 30分钟超时
	opts = append(opts, asynq.Retention(2*time.Hour))  // 保留2小时
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	opts = append(opts, asynq.Queue(string(group)))
	_, err := do.MustInvoke[*asynq.Scheduler](nil).Register(cron, task, opts...)
	if err != nil {
		panic("Scheduler add error :" + err.Error())
	}
}

// RegTask 注册队列任务
func RegTask(pattern string, handler func(context.Context, *asynq.Task) error) {
	do.MustInvoke[*asynq.ServeMux](nil).HandleFunc(pattern, handler)
}

type TaskLogger struct {
	Logger zerolog.Logger
}

func (t *TaskLogger) Debug(args ...interface{}) {
	t.Logger.Debug().Msg(fmt.Sprint(args...))
}

func (t *TaskLogger) Info(args ...interface{}) {
	t.Logger.Info().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Warn(args ...interface{}) {
	t.Logger.Warn().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Error(args ...interface{}) {
	t.Logger.Error().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Fatal(args ...interface{}) {
	t.Logger.Fatal().Msg(fmt.Sprint(args...))
}
