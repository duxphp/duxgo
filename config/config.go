package config

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

type Config map[string]*viper.Viper

func Init() {
	// 注册di服务
	do.ProvideValue[Config](nil, map[string]*viper.Viper{})

	pwd, _ := os.Getwd()
	configFiles, err := filepath.Glob(filepath.Join(pwd, registry.ConfigDir+"*.yaml"))
	if err != nil {
		panic("configuration loading failure")
	}

	// 循环加载配置文件
	for _, file := range configFiles {
		filename := path.Base(file)
		suffix := path.Ext(file)
		name := filename[0 : len(filename)-len(suffix)]
		do.MustInvoke[Config](nil)[name] = LoadConfig(name)
	}

	// 设置框架配置
	registry.Debug = registry.Config["app"].GetBool("server.debug")
	registry.DebugMsg = registry.Config["app"].GetString("server.debugMsg")

}

// LoadConfig 加载配置
func LoadConfig(name string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(name)
	config.SetConfigType("yaml")
	config.AddConfigPath(registry.ConfigDir)
	if err := config.ReadInConfig(); err != nil {
		fmt.Println("config", name)
		panic(err)
	}
	return config
}

// Get 获取配置
func Get(name string) *viper.Viper {
	if t, ok := do.MustInvoke[Config](nil)[name]; ok {
		return t
	} else {
		panic("configuration (" + name + ") not found")
	}
}
