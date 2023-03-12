package logger

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/config"
	"github.com/duxphp/duxgo/v2/helper"
	"github.com/rs/zerolog"
	"github.com/samber/do"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

// Log 日志
func Log() *zerolog.Logger {
	return do.MustInvoke[*zerolog.Logger](nil)
}

type loggerConfig struct {
	Path       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// Init 初始化日志
func Init() {
	path := config.Get("app").GetString("logger.default.path")
	if !helper.IsExist(path) {
		if !helper.CreateDir(path) {
			panic("failed to create log directory")
		}
	}

	// 默认日志配置
	logConfig := loggerConfig{
		Path:       path,
		MaxSize:    config.Get("app").GetInt("logger.default.maxSize"),
		MaxBackups: config.Get("app").GetInt("logger.default.maxBackups"),
		MaxAge:     config.Get("app").GetInt("logger.default.maxAge"),
		Compress:   config.Get("app").GetBool("logger.default.compress"),
	}

	// 初始化默认日志，根据日志等级分别输出
	writerList := make([]io.Writer, 0)
	levels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	for _, level := range levels {
		writerList = append(writerList, GetWriter(
			level,
			fmt.Sprintf("%s/%s.log", logConfig.Path, level),
			logConfig.MaxSize,
			logConfig.MaxBackups,
			logConfig.MaxAge,
			logConfig.Compress,
			false,
		))
	}
	log := New(writerList...).With().Timestamp().Caller().Logger()
	do.ProvideValue[*zerolog.Logger](nil, &log)

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
