package config

import (
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Server       *server // 服务器配置的全局变量
	Mysql        *mySQL  // MySQL数据库配置的全局变量
	MaxBorrowNum *maxBorrowNum
	runtimeViper *viper.Viper // Viper实例，用于管理配置文件
)

// Init 初始化配置模块
func Init() {
	runtimeViper = viper.New()             // 创建一个新的Viper实例
	configPath := "./config/config.yaml"   // 配置文件路径
	runtimeViper.SetConfigFile(configPath) // 设置配置文件路径
	runtimeViper.SetConfigType("yaml")     // 设置配置文件类型为YAML

	// 检查配置文件是否存在，如果不存在则创建默认配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			hlog.Errorf("config.Init: failed to create default config: %v", err) // 记录错误日志
			return
		}
		hlog.Info("config.Init: default config file created") // 记录信息日志
	}

	// 读取配置文件
	if err := runtimeViper.ReadInConfig(); err != nil {
		hlog.Infof("config.Init: config: read error: %v\n", err) // 记录错误日志
		return
	}
	configMapping() // 将配置映射到全局变量

	// 监听配置文件的变化
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		hlog.Infof("config: notice config changed: %v\n", e.String()) // 记录配置文件变化
		configMapping()                                               // 重新映射配置到全局变量
	})
	runtimeViper.WatchConfig() // 开始监听配置文件
}

// createDefaultConfig 创建默认的配置文件
func createDefaultConfig(configPath string) error {
	// 创建配置文件所在的目录
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	// 默认配置信息
	defaultConfig := config{
		Server: server{
			Addr: "127.0.0.1", // 默认服务器地址
			Port: 8080,        // 默认服务器端口
		},
		MySQL: mySQL{
			Addr:     "127.0.0.1:3306", // 默认数据库地址
			Database: "LMS",            // 默认数据库名称
			Username: "root",           // 默认数据库用户名
			Password: "root",           // 默认数据库密码
			Charset:  "utf8mb4",        // 默认数据库字符集
		},
		maxBorrowNum: maxBorrowNum{
			Num: 5,
		},
	}

	// 使用Viper将默认配置写入文件
	v := viper.New()
	v.Set("server", defaultConfig.Server)
	v.Set("mysql", defaultConfig.MySQL)
	v.Set("maxBorrowNum", defaultConfig.maxBorrowNum.MaxBorrowNum)

	return v.WriteConfigAs(configPath) // 将默认配置写入指定路径
}

// configMapping 将配置文件的内容映射到全局变量
func configMapping() {
	c := new(config) // 创建一个新的配置对象
	// 将Viper中的配置解码到配置对象中
	if err := runtimeViper.Unmarshal(&c); err != nil {
		hlog.Infof("config.configMapping: config: unmarshal error: %v", err) // 记录错误日志
	}
	// 将配置对象的值赋给全局变量
	Server = &c.Server
	Mysql = &c.MySQL
	MaxBorrowNum = &c.MaxBorrowNum
}
