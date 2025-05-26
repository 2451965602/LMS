package config

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	Server       *server
	Mysql        *mySQL
	runtimeViper *viper.Viper
)

func Init() {
	runtimeViper = viper.New()
	configPath := "./config/config.yaml"
	runtimeViper.SetConfigFile(configPath)
	runtimeViper.SetConfigType("yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			hlog.Errorf("config.Init: failed to create default config: %v", err)
			return
		}
		hlog.Info("config.Init: default config file created")
		os.Exit(0)
	}

	if err := runtimeViper.ReadInConfig(); err != nil {
		hlog.Infof("config.Init: config: read error: %v\n", err)
		return
	}
	configMapping()

	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		hlog.Infof("config: notice config changed: %v\n", e.String())
		configMapping() // 重新映射配置
	})
	runtimeViper.WatchConfig()
}

func createDefaultConfig(configPath string) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	defaultConfig := config{
		Server: server{
			Addr: "127.0.0.1",
			Port: 8080,
		},
		MySQL: mySQL{
			Addr:     "127.0.0.1:3306",
			Database: "LMS",
			Username: "root",
			Password: "root",
			Charset:  "utf8mb4",
		},
	}

	v := viper.New()
	v.Set("server", defaultConfig.Server)
	v.Set("mysql", defaultConfig.MySQL)

	return v.WriteConfigAs(configPath)
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
