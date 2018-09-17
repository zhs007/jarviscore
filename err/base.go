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

func (err *jarvisError) Error() string {
	return "JarvisError - [" + string(err.errcode) + "]"
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

// const (
// 	FILEREADSIZEINVALID = 1
// 	PEERADDREMPTY       = 2
// 	NOINIT              = 3
// )

// // GetErrCodeString -
// func GetErrCodeString(errcode int) string {
// 	switch errcode {
// 	case FILEREADSIZEINVALID:
// 		return "invalid filesize & bytesread."
// 	case PEERADDREMPTY:
// 		return "peeraddr is empty."
// 	case NOINIT:
// 		return "no init."
// 	}

// 	return ""
// }
