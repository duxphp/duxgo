package config

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

func Init() {

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
		registry.Config[name] = LoadConfig(name)
	}

	// 设置框架配置
	registry.Debug = registry.Config["app"].GetBool("server.debug")
	registry.DebugMsg = registry.Config["app"].GetString("server.debugMsg")

}

// LoadConfig 加载配置
func LoadConfig(name string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(name)
	config.SetConfigType("toml")
	config.AddConfigPath(registry.ConfigDir)
	if err := config.ReadInConfig(); err != nil {
		fmt.Println("config", name)
		panic(err)
	}
	return config
}

// Get 获取配置
func Get(name string) *viper.Viper {
	if t, ok := registry.Config[name]; ok {
		return t
	} else {
		panic("configuration (" + name + ") not found")
	}
}
