package jarviscore

import (
	"os/exec"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	// CtrlTypeScriptFile - scriptfile ctrltype
	CtrlTypeScriptFile = "scriptfile"
)

// CtrlScriptFile -
type CtrlScriptFile struct {
}

// runScript
func (ctrl *CtrlScriptFile) runScript(ci *pb.CtrlInfo) ([]byte, error) {
	var csd pb.CtrlScriptData
	err := ptypes.UnmarshalAny(ci.Dat, &csd)
	if err != nil {
		return nil, err
	}

	return exec.Command("sh", "-c", string(csd.File)).CombinedOutput()
}

// Run -
func (ctrl *CtrlScriptFile) Run(jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {

	var msgs []*pb.JarvisMsg

	out, err := ctrl.runScript(ci)
	if err != nil {
		return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, AppendString(string(out), err.Error()), msgs)
	}

	return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, string(out), msgs)
}

// BuildCtrlInfoForScriptFile - build ctrlinfo for scriptfile
// Deprecated: you can use BuildCtrlInfoForScriptFile3
func BuildCtrlInfoForScriptFile(filename string, filedata []byte,
	destpath string) (*pb.CtrlInfo, error) {

	csd := &pb.CtrlScriptData{
		File:     filedata,
		DestPath: destpath,
		Filename: filename,
	}

	dat, err := ptypes.MarshalAny(csd)
	if err != nil {
		return nil, err
	}

	ci := &pb.CtrlInfo{
		CtrlType: CtrlTypeScriptFile,
		Dat:      dat,
	}

	return ci, nil
}
