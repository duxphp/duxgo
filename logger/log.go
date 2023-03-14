package logger

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/config"
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

// Init 初始化日志
func Init() {
	// 初始化默认日志，根据日志等级分别输出
	writerList := make([]io.Writer, 0)
	levels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	for _, level := range levels {
		writerList = append(writerList, GetWriter(
			level,
			"default",
			level,
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
func GetWriter(level string, dirName string, name string, recursion bool) *LevelWriter {
	parseLevel, _ := zerolog.ParseLevel(level)
	return &LevelWriter{zerolog.MultiLevelWriter(&lumberjack.Logger{
		Filename:   fmt.Sprintf("./data/%s/%s.log", dirName, name),        // 日志文件路径
		MaxSize:    config.Get("app").GetInt("logger.default.maxSize"),    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.Get("app").GetInt("logger.default.maxBackups"), // 日志文件最多保存多少个备份
		MaxAge:     config.Get("app").GetInt("logger.default.maxAge"),     // 文件最多保存多少天
		Compress:   config.Get("app").GetBool("logger.default.compress"),  // 是否压缩
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
