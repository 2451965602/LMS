package config

// Server 用于存储服务器配置信息
type server struct {
	Addr string `yaml:"addr"` // 服务器地址
	Port int64  `yaml:"port"` // 服务器端口
}

// mySQL 用于存储MySQL数据库配置信息
type mySQL struct {
	Addr     string `yaml:"addr"`     // 数据库地址
	Database string `yaml:"database"` // 数据库名称
	Username string `yaml:"username"` // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	Charset  string `yaml:"charset"`  // 数据库字符集
}
type maxBorrowNum struct {
	Num int64 `yaml:"Num"`
}

// config 用于存储整个配置信息
type config struct {
	Server       server       `yaml:"server"` // 服务器配置
	MySQL        mySQL        `yaml:"mysql"`  // MySQL数据库配置
	maxBorrowNum maxBorrowNum `yaml:"maxBorrowNum"`
}
