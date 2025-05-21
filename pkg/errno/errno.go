package errno

import (
	"errors"
	"fmt"
)

type ErrNo struct {
	ErrorCode int64
	ErrorMsg  string
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

func (e ErrNo) Error() string { return e.ErrorMsg }

func NewErrNoWithStack(code int64, msg string) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

func Errorf(code int64, template string, args ...interface{}) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  fmt.Sprintf(template, args...),
	}
}

func (e ErrNo) WithMessage(message string) ErrNo {
	e.ErrorMsg = message
	return e
}

func (e ErrNo) WithError(err error) ErrNo {
	e.ErrorMsg = err.Error()
	return e
}

func ConvertErr(err error) ErrNo {
	if err == nil {
		return Success
	}
	errno := ErrNo{}
	if errors.As(err, &errno) {
		return errno
	}
	s := InternalServiceError
	s.ErrorMsg = err.Error()
	return s
}
