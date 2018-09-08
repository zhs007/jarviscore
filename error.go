package jarviscore

import (
	"github.com/zhs007/jarviscore/errcode"
	"github.com/zhs007/jarviscore/log"
	"go.uber.org/zap"
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
	err := &jarvisError{errcode: errcode}

	log.Error(err.Error(), zap.Int("ErrCode", errcode))

	return err
}

// IsJarvisError -
func IsJarvisError(err interface{}) bool {
	if _, ok := err.(JarvisError); ok {
		return true
	}

	return false
}

// ErrorLog -
func ErrorLog(msg string, err error) {
	if jerr, ok := err.(JarvisError); ok {
		log.Error(msg, zap.String("err", jerr.Error()), zap.Int("errcode", jerr.ErrCode()))
	} else {
		log.Error(msg, zap.String("err", err.Error()))
	}
}

// warnLog -
func warnLog(msg string, err error) {
	if jerr, ok := err.(JarvisError); ok {
		log.Warn(msg, zap.String("err", jerr.Error()), zap.Int("errcode", jerr.ErrCode()))
	} else {
		log.Warn(msg, zap.String("err", err.Error()))
	}
}
