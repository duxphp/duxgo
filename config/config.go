package config

import (
	"fmt"
	"github.com/duxphp/duxgo/v2/core"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

func Init() {

	// 解析配置文件
	pwd, _ := os.Getwd()
	configFiles, err := filepath.Glob(filepath.Join(pwd, core.ConfigDir+"*.toml"))
	if err != nil {
		panic("configuration loading failure")
	}

	for _, file := range configFiles {
		filename := path.Base(file)
		suffix := path.Ext(file)
		name := filename[0 : len(filename)-len(suffix)]
		core.Config[name] = LoadConfig(name)
	}

	// 调试配置
	core.Debug = core.Config["app"].GetBool("server.debug")
	core.DebugMsg = core.Config["app"].GetString("server.debugMsg")

}

// LoadConfig 加载配置
func LoadConfig(name string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(name)
	config.SetConfigType("toml")
	config.AddConfigPath(core.ConfigDir)
	if err := config.ReadInConfig(); err != nil {
		fmt.Println("config", name)
		panic(err)
	}
	return config
}

// Get 获取配置
func Get(name string) *viper.Viper {
	if t, ok := core.Config[name]; ok {
		return t
	} else {
		panic("configuration (" + name + ") not found")
	}
}
