package jarviscore

import (
	"os/exec"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	// CtrlTypeScriptFile3 - scriptfile3 ctrltype
	CtrlTypeScriptFile3 = "scriptfile3"
)

// CtrlScriptFile3 -
type CtrlScriptFile3 struct {
}

// runScript
func (ctrl *CtrlScriptFile3) runScript(ci *pb.CtrlInfo) ([]byte, error) {
	var csd3 pb.CtrlScript3Data
	err := ptypes.UnmarshalAny(ci.Dat, &csd3)
	if err != nil {
		return nil, err
	}

	return exec.Command("sh", "-c", string(csd3.ScriptFile.File)).CombinedOutput()
}

// Run -
func (ctrl *CtrlScriptFile3) Run(jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {

	out, err := ctrl.runScript(ci)
	if err != nil {
		return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, AppendString(string(out), err.Error()))
	}

	return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, string(out))
}

// BuildCtrlInfoForScriptFile3 - build ctrlinfo for scriptfile
func BuildCtrlInfoForScriptFile3(scriptfile *pb.FileData, endFiles []string) (*pb.CtrlInfo, error) {

	csd3 := &pb.CtrlScript3Data{
		ScriptFile: scriptfile,
		EndFiles:   endFiles,
	}

	dat, err := ptypes.MarshalAny(csd3)
	if err != nil {
		return nil, err
	}

	ci := &pb.CtrlInfo{
		CtrlType: CtrlTypeScriptFile3,
		Dat:      dat,
	}

	return ci, nil
}
