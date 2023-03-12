package cache

import (
	"github.com/coocood/freecache"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/samber/do"
	"runtime/debug"
)

func Init() {
	// 缓存大小，单位 M
	cacheSize := registry.Config["app"].GetInt("cache.size") * 1024 * 1024
	// 注册di服务
	do.ProvideValue[*freecache.Cache](nil, freecache.NewCache(cacheSize))
	debug.SetGCPercent(20)

}
