package task

import (
	"context"
	"github.com/duxphp/duxgo/core"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/cast"
	"time"
)

func Init() {
	dbConfig := core.Config["database"].GetStringMapString("redis")
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
			Logger:      &TaskLogger{},
			LogLevel:    asynq.WarnLevel,
			Concurrency: 20,
			Queues: map[string]int{
				"high":    10,
				"default": 7,
				"low":     3,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				retried, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				if retried >= maxRetry {
					core.Logger.Info().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("task retry")
				} else {
					core.Logger.Info().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("task error")
				}
			}),
		},
	)

	// 混合器
	mux := asynq.NewServeMux()

	// 中间件
	//mux.Use(func(next asynq.Handler) asynq.Handler {
	//	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
	//		return nil
	//	})
	//})

	// 客户端
	client := asynq.NewClient(res)

	// 检查器
	inspector := asynq.NewInspector(res)

	core.QueueMux = mux
	core.Queue = srv
	core.QueueClient = client
	core.QueueInspector = inspector

	core.QueueMux.HandleFunc("ping", func(ctx context.Context, t *asynq.Task) error {
		core.Logger.Debug().Msg("queue ping status")
		return nil
	})

	// 定时调度服务
	scheduler := asynq.NewScheduler(res, &asynq.SchedulerOpts{
		LogLevel: asynq.ErrorLevel,
		Location: core.TimeLocation,
		EnqueueErrorHandler: func(task *asynq.Task, opts []asynq.Option, err error) {
			core.Logger.Error().Msgf("scheduler: ", err.Error())
		},
	})
	core.Scheduler = scheduler

}

type Priority string

const (
	PRIORITY_HIGH    Priority = "high"
	PRIORITY_DEFAULT Priority = "default"
	PRIORITY_LOW     Priority = "low"
)

// StartQueue 启动队列服务
func StartQueue() {
	if err := core.Queue.Run(core.QueueMux); err != nil {
		core.Logger.Error().Msgf("Queue service cannot be started: %v", err)
	}
}

// StartScheduler 启动调度服务
func StartScheduler() {
	if err := core.Scheduler.Run(); err != nil {
		core.Logger.Error().Msgf("Scheduler service cannot be started: %v", err)
	}
}

// Add 添加即时队列
func Add(typename string, params any, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return addTask(typename, params, asynq.Queue(string(group)))
}

// AddDelay 延迟队列（秒）
func AddDelay(typename string, params any, t time.Duration, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return addTask(typename, params, asynq.ProcessIn(t), asynq.Queue(string(group)))
}

// AddTime 定时队列（秒）
func AddTime(typename string, params any, t time.Time, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return addTask(typename, params, asynq.ProcessAt(t), asynq.Queue(string(group)))
}

// 添加队列任务
func addTask(typename string, params any, opts ...asynq.Option) *asynq.TaskInfo {
	payload, _ := json.Marshal(params)
	task := asynq.NewTask(typename, payload)
	opts = append(opts, asynq.MaxRetry(3))            // 重试3次
	opts = append(opts, asynq.Timeout(1*time.Minute)) // 1分钟超时
	opts = append(opts, asynq.Retention(2*time.Hour)) // 保留2小时

	info, err := core.QueueClient.Enqueue(task, opts...)
	if err != nil {
		core.Logger.Error().Msg("Queue add error :" + err.Error())
	}
	return info
}

// DelTask 删除队列任务
func DelTask(priority Priority, id string) error {
	err := core.QueueInspector.DeleteTask(string(priority), id)
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
	_, err := core.Scheduler.Register(cron, task, opts...)
	if err != nil {
		panic("Scheduler add error :" + err.Error())
	}
}

type TaskLogger struct {
}

func (t *TaskLogger) Debug(args ...interface{}) {
	core.Logger.Debug().Msg(fmt.Sprint(args...))
}

func (t *TaskLogger) Info(args ...interface{}) {
	core.Logger.Info().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Warn(args ...interface{}) {
	core.Logger.Warn().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Error(args ...interface{}) {
	core.Logger.Error().Interface("args", args).Msg("task")

}

func (t *TaskLogger) Fatal(args ...interface{}) {
	core.Logger.Fatal().Msg(fmt.Sprint(args...))

}
