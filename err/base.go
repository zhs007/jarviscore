package jarviserr

import (
	"github.com/zhs007/jarviscore/log"
	pb "github.com/zhs007/jarviscore/proto"
	"go.uber.org/zap"
)

// JarvisError -
type JarvisError interface {
	Error() string
	ErrCode() pb.CODE
}

// jarvisError -
type jarvisError struct {
	errcode pb.CODE
}

// FormatCode -
func FormatCode(code pb.CODE) pb.CODE {
	if _, ok := pb.CODE_name[int32(code)]; !ok {
		return pb.CODE_INVALID_CODE
	}

	return code
}

func (err *jarvisError) Error() string {
	return "JarvisError - [" + pb.CODE_name[int32(FormatCode(err.errcode))] + "]"
}

func (err *jarvisError) ErrCode() pb.CODE {
	return err.errcode
}

// NewError -
func NewError(errcode pb.CODE) JarvisError {
	err := &jarvisError{errcode: errcode}

	log.Error(err.Error(), zap.Int("ErrCode", int(errcode)))

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
		log.Error(msg, zap.String("err", jerr.Error()), zap.Int("errcode", int(jerr.ErrCode())))
	} else {
		log.Error(msg, zap.String("err", err.Error()))
	}
}

// WarnLog -
func WarnLog(msg string, err error) {
	if jerr, ok := err.(JarvisError); ok {
		log.Warn(msg, zap.String("err", jerr.Error()), zap.Int("errcode", int(jerr.ErrCode())))
	} else {
		log.Warn(msg, zap.String("err", err.Error()))
	}
}
