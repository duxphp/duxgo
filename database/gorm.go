package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/duxphp/duxgo/core"
	coreLogger "github.com/duxphp/duxgo/logger"
	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func GormInit() {
	dbConfig := core.Config["database"].GetStringMapString("db")

	var connect gorm.Dialector
	if dbConfig["type"] == "mysql" {
		connect = mysql.Open(dbConfig["username"] + ":" + dbConfig["password"] + "@tcp(" + dbConfig["host"] + ":" + dbConfig["port"] + ")/" + dbConfig["dbname"] + "?charset=utf8mb4&parseTime=True&loc=Local")
	}
	if dbConfig["type"] == "postgresql" {
		connect = postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbConfig["host"],
			dbConfig["username"],
			dbConfig["password"],
			dbConfig["dbname"],
			dbConfig["port"],
		))
	}
	if dbConfig["type"] == "sqlite" {
		connect = sqlite.Open(dbConfig["file"])
	}
	database, err := gorm.Open(connect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "app_",
			SingularTable: true,
		},
		Logger: GormLogger(),
	})
	if err != nil {
		panic("database error: " + err.Error())
	}
	core.Db = database
	sqlDB, err := core.Db.DB()
	if err != nil {
		panic("database error: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(core.Config["app"].GetInt("database.maxIdleConns"))
	sqlDB.SetMaxOpenConns(core.Config["app"].GetInt("database.maxOpenConns"))

}

type logger struct {
	SlowThreshold             time.Duration
	SourceField               string
	IgnoreRecordNotFoundError bool
	Logger                    zerolog.Logger
	LogLevel                  gormLogger.LogLevel
}

func GormLogger() *logger {
	vLog := coreLogger.New(
		core.Config["app"].GetString("logger.db.level"),
		core.Config["app"].GetString("logger.db.path"),
		core.Config["app"].GetInt("logger.db.maxSize"),
		core.Config["app"].GetInt("logger.db.maxBackups"),
		core.Config["app"].GetInt("logger.db.maxAge"),
		core.Config["app"].GetBool("logger.db.compress"),
	).With().Caller().CallerWithSkipFrameCount(5).Timestamp().Logger()

	return &logger{
		SlowThreshold:             1 * time.Second,
		Logger:                    vLog,
		LogLevel:                  gormLogger.Silent,
		IgnoreRecordNotFoundError: true,
	}
}

func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &logger{
		Logger:                    l.Logger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel < gormLogger.Info {
		return
	}
	l.Logger.Info().Msgf(s, args)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel < gormLogger.Warn {
		return
	}
	l.Logger.Warn().Msgf(s, args)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel < gormLogger.Error {
		return
	}
	l.Logger.Error().Msgf(s, args)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := map[string]interface{}{
		"sql":      sql,
		"duration": elapsed,
	}
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.Logger.Error().Err(err).Fields(fields).Msg("[GORM] query error")
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		l.Logger.Warn().Fields(fields).Msgf("[GORM] slow query")
	case l.LogLevel >= gormLogger.Info:
		l.Logger.Debug().Fields(fields).Msgf("[GORM] query")
	}

}
