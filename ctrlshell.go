package jarviscore

import (
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	// CtrlTypeShell - shell ctrltype
	CtrlTypeShell = "shell"
)

// CtrlShell -
type CtrlShell struct {
}

// Run -
func (ctrl *CtrlShell) Run(ci *pb.CtrlInfo) ([]byte, error) {
	// cmd := exec.Command("whoami")
	// whoami, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
