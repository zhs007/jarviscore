package jarviscore

import (
	"context"
	"os/exec"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/zhs007/jarviscore/proto"
)

const (
	// CtrlTypeScriptFile2 - scriptfile2 ctrltype
	CtrlTypeScriptFile2 = "scriptfile2"
)

// CtrlScriptFile2 -
type CtrlScriptFile2 struct {
}

// runScript
func (ctrl *CtrlScriptFile2) runScript(ci *pb.CtrlInfo) ([]byte, error) {
	var csd2 pb.CtrlScript2Data
	err := ptypes.UnmarshalAny(ci.Dat, &csd2)
	if err != nil {
		return nil, err
	}

	for _, v := range csd2.SrcFiles {
		err = StoreLocalFile(v)
		if err != nil {
			return nil, err
		}
	}

	return exec.Command("sh", "-c", string(csd2.ScriptFile.File)).CombinedOutput()
}

// Run -
func (ctrl *CtrlScriptFile2) Run(ctx context.Context, jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {

	var msgs []*pb.JarvisMsg

	out, err := ctrl.runScript(ci)
	if err != nil {
		return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, AppendString(string(out), err.Error()), msgs)
	}

	return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, string(out), msgs)
}

// BuildCtrlInfoForScriptFile2 - build ctrlinfo for scriptfile
// Deprecated: you can use BuildCtrlInfoForScriptFile3
func BuildCtrlInfoForScriptFile2(scriptfile *pb.FileData, files []*pb.FileData) (*pb.CtrlInfo, error) {

	csd2 := &pb.CtrlScript2Data{
		ScriptFile: scriptfile,
		SrcFiles:   files,
	}

	dat, err := ptypes.MarshalAny(csd2)
	if err != nil {
		return nil, err
	}

	ci := &pb.CtrlInfo{
		CtrlType: CtrlTypeScriptFile2,
		Dat:      dat,
	}

	return ci, nil
}
