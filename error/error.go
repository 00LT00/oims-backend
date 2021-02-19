package error

import (
	"fmt"
	"runtime"
)

type Error struct {
	ErrCode string
	Msg     string
}

func caller() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误是: %v, 发生在: %s", e, caller())
}

func NewError(ErrCode, Msg string) Error {
	e := Error{ErrCode: ErrCode, Msg: Msg}
	return e
}
