package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxphp/duxgo/global"
	"runtime/debug"
)

func Init() {
	// 缓存大小，单位 M
	cacheSize := global.Config["app"].GetInt("cache.size") * 1024 * 1024
	global.Cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)

}
