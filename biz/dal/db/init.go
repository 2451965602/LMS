package db

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/2451965602/LMS/pkg/constants"
	"github.com/2451965602/LMS/pkg/errno"
	"github.com/2451965602/LMS/pkg/utils"
)

var db *gorm.DB

func Init() error {
	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL get mysql DSN error: %v", err))
	}

	db, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,  // 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
			SkipDefaultTransaction: false, // 不禁用默认事务(即单个创建、更新、删除时使用事务)
			TranslateError:         true,  // 允许翻译错误
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.InitMySQL mysql connect error: %v", err)
	}

	sqlDB, err := db.DB() // 尝试获取 DB 实例对象
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, fmt.Sprintf("get generic database object error: %v", err))
	}

	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)       // 最大闲置连接数
	sqlDB.SetMaxOpenConns(constants.MaxConnections)     // 最大连接数
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime) // 最大可复用时间
	sqlDB.SetConnMaxIdleTime(constants.ConnMaxIdleTime) // 最长保持空闲状态时间
	db = db.WithContext(context.Background())

	// 进行连通性测试
	if err = sqlDB.Ping(); err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, fmt.Sprintf("ping database error: %v", err))
	}
	return nil
}
