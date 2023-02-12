package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxphp/duxgo/v2/registry"
	"runtime/debug"
)

func Init() {
	// 缓存大小，单位 M
	cacheSize := registry.Config["app"].GetInt("cache.size") * 1024 * 1024
	registry.Cache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)

}
