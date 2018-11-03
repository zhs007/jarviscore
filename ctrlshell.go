package jarviscore

import (
	"os/exec"
)

// CtrlShell -
type CtrlShell struct {
}

// Run -
func (ctrl *CtrlShell) Run(command []byte) ([]byte, error) {
	cmd := exec.Command("whoami")
	whoami, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return whoami, nil
}
