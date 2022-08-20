package logger

import (
	"fmt"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/util/function"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type loggerConfig struct {
	Path       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      string
}

// Init 初始化日志
func Init() {
	path := core.Config["app"].GetString("logger.default.path")
	if !function.IsExist(path) {
		if !function.CreateDir(path) {
			panic("failed to create log directory")
		}
	}

	// 默认日志配置
	config := loggerConfig{
		Path:       path,
		MaxSize:    core.Config["app"].GetInt("logger.default.maxSize"),
		MaxBackups: core.Config["app"].GetInt("logger.default.maxBackups"),
		MaxAge:     core.Config["app"].GetInt("logger.default.maxAge"),
		Compress:   core.Config["app"].GetBool("logger.default.compress"),
		Level:      core.Config["app"].GetString("logger.default.level"),
	}
	core.Logger = New(GetWriter(
		config.Level,
		fmt.Sprintf("%s/app.log", config.Path),
		config.MaxSize,
		config.MaxBackups,
		config.MaxAge,
		config.Compress,
		true,
	)).With().Timestamp().Caller().Logger()
}

// New 新建日志
func New(writers ...io.Writer) zerolog.Logger {
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	writers = append(writers, &console)
	multi := zerolog.MultiLevelWriter(writers...)
	return zerolog.New(multi)
}

// GetWriter 获取日志驱动
func GetWriter(level string, path string, maxSize int, maxBackups int, maxAge int, compress bool, recursion bool) *LevelWriter {
	parseLevel, _ := zerolog.ParseLevel(level)
	return &LevelWriter{zerolog.MultiLevelWriter(&lumberjack.Logger{
		Filename:   path,       // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}), parseLevel, recursion}
}

type LevelWriter struct {
	w         zerolog.LevelWriter
	level     zerolog.Level
	recursion bool
}

func (w *LevelWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}
func (w *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.level && w.recursion {
		return w.w.WriteLevel(level, p)
	}
	if level == w.level && !w.recursion {
		return w.w.WriteLevel(level, p)
	}
	return len(p), nil
}
