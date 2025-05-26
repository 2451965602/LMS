package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/2451965602/LMS/config"
)

func GetMysqlDSN() (string, error) {
	if config.Mysql == nil {
		return "", errors.New("config not found")
	}

	dsn := strings.Join([]string{
		config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/",
		config.Mysql.Database, "?charset=" + config.Mysql.Charset + "&parseTime=true",
	}, "")

	return dsn, nil
}

func GetServerAddr() (string, error) {
	if config.Server == nil {
		return "", errors.New("config not found")
	}

	addr := strings.Join([]string{
		config.Server.Addr, ":", strconv.FormatInt(config.Server.Port, 10),
	}, "")

	return addr, nil
}
