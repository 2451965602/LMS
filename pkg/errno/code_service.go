package errno

// 业务强相关, 范围是 1000-9999
const (
	ServiceUserExist = 1000 + iota
	ServiceUserNotExist
	ServiceBookTypeNotFound
	ServiceBookNotExist
	ServiceBorrowRecordNotExist
	ServicePermissionDenied
	ServiceBookTypeExist
	ServiceBookTypeNotExist
	ServiceBookNotAvailable
	ServiceReserveDateError
	ServiceActionNotAllowed
	ServiceBookTypeInUse
)
