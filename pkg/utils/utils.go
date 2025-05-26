package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/2451965602/LMS/config"
)

// GetMysqlDSN 生成MySQL数据库的DSN（Data Source Name）字符串
// 返回值：
//   - string: MySQL数据库的DSN字符串
//   - error: 错误信息，如果配置未找到会返回错误
func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("config not found") // 如果MySQL配置未找到，返回错误
	}

	// 拼接MySQL DSN字符串
	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password, // 用户名和密码
		"@tcp(", config.Mysql.Addr, ")/", // 数据库地址
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true", // 数据库名称、字符集和解析时间
	}, "")

	return dsn, nil
}

// GetServerAddr 生成服务器的地址字符串
// 返回值：
//   - string: 服务器的地址字符串
//   - error: 错误信息，如果配置未找到会返回错误
func GetServerAddr() (string, error) {
	if config.Server == nil {
		return "", errors.New("config not found") // 如果服务器配置未找到，返回错误
	}

	// 拼接服务器地址字符串
	addr := strings.Join([]string{
		config.Server.Addr, ":", strconv.FormatInt(config.Server.Port, 10), // 服务器地址和端口
	}, "")

	return addr, nil
}
