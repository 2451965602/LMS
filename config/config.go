package config

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Server       *server
	Mysql        *mySQL
	runtimeViper *viper.Viper
)

func Init() {
	runtimeViper = viper.New()
	runtimeViper.SetConfigFile("./config/config.yaml")
	runtimeViper.SetConfigType("yaml")

	if err := runtimeViper.ReadInConfig(); err != nil {
		hlog.Infof("config.Init: config: read error: %v\n", err)
		return
	}
	configMapping()

	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		// 我们无法确定监听到配置变更时是否已经初始化完毕，所以此处需要做一个判断
		hlog.Infof("config: notice config changed: %v\n", e.String())
		configMapping() // 重新映射配置
	})
	runtimeViper.WatchConfig()
}

// configMapping 用于将配置映射到全局变量
func configMapping() {
	c := new(config)
	if err := runtimeViper.Unmarshal(&c); err != nil {
		// 由于这个函数会在配置重载时被再次触发，所以需要判断日志记录方式
		hlog.Infof("config.configMapping: config: unmarshal error: %v", err)
	}
	Server = &c.Server
	Mysql = &c.MySQL
}
