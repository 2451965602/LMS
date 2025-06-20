package errno

// 业务强相关, 范围是 1000-9999
const (
	ServiceUserExist = 1000 + iota
	ServiceUserNotExist
	ServiceInvalidUsername
	ServiceInvalidPhone
	ServicePermissionDenied

	ServiceInvalidISBN
	ServiceInvalidAuthor
	ServiceBookTypeNotFound
	ServiceBookTypeInUse
	ServiceBookTypeExist
	ServiceBookTypeNotExist

	ServiceBookNotExist
	ServiceBookNotAvailable

	ServiceBorrowRecordNotExist

	ServiceActionNotAllowed
)
