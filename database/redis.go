package database

import (
	"github.com/duxphp/duxgo/global"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

func RedisInit() {
	dbConfig := global.Config["database"].GetStringMapString("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       cast.ToInt(dbConfig["db"]),
	})
	_, err := client.Ping(global.Ctx).Result()
	if err != nil {
		panic(err.Error())
	}

	global.Redis = client
}
