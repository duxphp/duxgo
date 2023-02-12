package database

import (
	"github.com/duxphp/duxgo/core"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

func RedisInit() {
	dbConfig := core.Config["database"].GetStringMapString("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       cast.ToInt(dbConfig["db"]),
	})
	_, err := client.Ping(core.Ctx).Result()
	if err != nil {
		panic(err.Error())
	}

	core.Redis = client
}
