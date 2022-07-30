package logger

import (
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/util/function"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func Init() {
	path := core.Config["app"].GetString("logger.default.path")
	if !function.IsExist(path) {
		if !function.CreateDir(path) {
			panic("failed to create log directory")
		}
	}
	core.Logger = New(
		core.Config["app"].GetString("logger.default.level"),
		path,
		core.Config["app"].GetInt("logger.default.maxSize"),
		core.Config["app"].GetInt("logger.default.maxBackups"),
		core.Config["app"].GetInt("logger.default.maxAge"),
		core.Config["app"].GetBool("logger.default.compress"),
	).With().Timestamp().Caller().Logger()
}

func New(level string, path string, maxSize int, maxBackups int, maxAge int, compress bool) zerolog.Logger {
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	parseLevel, _ := zerolog.ParseLevel(level)
	fileLog := &LevelWriter{zerolog.MultiLevelWriter(&lumberjack.Logger{
		Filename:   path,       // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}), parseLevel}

	multi := zerolog.MultiLevelWriter(&console, fileLog)
	return zerolog.New(multi)
}

type LevelWriter struct {
	w     zerolog.LevelWriter
	level zerolog.Level
}

func (w *LevelWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}
func (w *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.level {
		return w.w.WriteLevel(level, p)
	}
	return len(p), nil
}
