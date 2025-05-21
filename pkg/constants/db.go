package constants

import "time"

const (
	MaxConnections  = 1000             // (DB) 最大连接数
	MaxIdleConns    = 10               // (DB) 最大空闲连接数
	ConnMaxLifetime = 10 * time.Second // (DB) 最大可复用时间
	ConnMaxIdleTime = 5 * time.Minute  // (DB) 最长保持空闲状态时间

	UserTableName         = "Users"         // (DB) 用户表名
	BookTypeTableName     = "BookTypes"     // (DB) 图书类型表名
	BookTableName         = "Books"         // (DB) 图书表名
	BorrowRecordTableName = "BorrowRecords" // (DB) 借阅记录表名
	ReservationTableName  = "Reservations"  // (DB) 预约记录表名

)
