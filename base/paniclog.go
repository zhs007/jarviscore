// +build !windows

package jarvisbase

import (
	"os"
	"path"
	"syscall"

	"go.uber.org/zap"
)

func initPanicFile() error {
	file, err := os.OpenFile(
		path.Join(logPath, BuildLogFilename("panic", logSubName)),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Warn("initPanicFile:OpenFile",
			zap.Error(err))

		return err
	}

	panicFile = file

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		Warn("initPanicFile:Dup2",
			zap.Error(err))

		return err
	}

	return nil
}
