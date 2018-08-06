package jarviscore

import (
	"github.com/zhs007/jarviscore/errcode"
)

// JarvisError -
type JarvisError interface {
	Error() string
	ErrCode() int
}

// jarvisError -
type jarvisError struct {
	errcode int
}

func (err *jarvisError) Error() string {
	return "JarvisError - [" + string(err.errcode) + "-" + jarviserrcode.GetErrCodeString(err.errcode) + "]"
}

func (err *jarvisError) ErrCode() int {
	return err.errcode
}

func newError(errcode int) JarvisError {
	return &jarvisError{errcode: errcode}
}

// IsJarvisError -
func IsJarvisError(err interface{}) bool {
	if _, ok := err.(JarvisError); ok {
		return true
	}

	return false
}
