package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxphp/duxgo/core"
	"runtime/debug"
)

func Init() {
	// 缓存大小，单位 M
	cacheSize := core.Config["app"].GetInt("cache.size") * 1024 * 1024
	core.Cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)

}
