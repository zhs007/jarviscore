package jarvisbase

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime/debug"
	"time"
)

// DumpPanic - dump panic file
func DumpPanic(err error) {
	ct := time.Now()

	filename := fmt.Sprintf("%v.panic", ct.String())

	f, err := os.Create(path.Join(logPath, filename))
	if err != nil {
		return
	}

	defer f.Close()

	io.WriteString(f, fmt.Sprintf("painc: %v\n", err.Error()))

	stack := debug.Stack()

	io.WriteString(f, fmt.Sprintf("stack: %v\n", string(stack)))
}
