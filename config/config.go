package config

import (
	"fmt"
	"github.com/duxphp/duxgo/core"
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
	fmt.Println("configFiles", configFiles)

	for _, file := range configFiles {
		filename := path.Base(file)
		suffix := path.Ext(file)
		name := filename[0 : len(filename)-len(suffix)]
		core.Config[name] = LoadConfig(name)
	}

	// 解析媒体文件
	//jsonPath, err := duxgo.StaticFs.Open("public/manifest.json")
	//if err != nil {
	//	panic(err.Error())
	//}
	//config := viper.New()
	//config.SetConfigType("json")
	//err = config.ReadConfig(jsonPath)
	//if err != nil {
	//	panic(err.Error())
	//}
	//jsonPath.Close()
	//duxgo.ConfigManifest = config.GetStringMap("src/main.js")

	// 调试配置
	core.Debug = core.Config["app"].GetBool("server.debug")
	core.DebugMsg = core.Config["app"].GetString("server.debugMsg")

}

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
