package jarviscore

import (
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

// Run -
func (ctrl *CtrlScriptFile2) Run(jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {
	var csd2 pb.CtrlScript2Data
	err := ptypes.UnmarshalAny(ci.Dat, &csd2)
	if err != nil {
		return BuildReply2ForCtrl(jarvisnode, srcAddr, msgid, pb.REPLYTYPE_ERROR, err.Error())
	}

	for _, v := range csd2.SrcFiles {
		err = StoreLocalFile(v)
		if err != nil {
			return BuildReply2ForCtrl(jarvisnode, srcAddr, msgid, pb.REPLYTYPE_ERROR, err.Error())
		}
	}

	out, err := exec.Command("sh", "-c", string(csd2.ScriptFile.File)).CombinedOutput()
	if err != nil {
		return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, AppendString(string(out), err.Error()))
	}

	return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, string(out))
}

// BuildCtrlInfoForScriptFile2 - build ctrlinfo for scriptfile
func BuildCtrlInfoForScriptFile2(ctrlid int64, scriptfile *pb.FileData, files []*pb.FileData) (*pb.CtrlInfo, error) {

	csd2 := &pb.CtrlScript2Data{
		ScriptFile: scriptfile,
		SrcFiles:   files,
	}

	dat, err := ptypes.MarshalAny(csd2)
	if err != nil {
		return nil, err
	}

	ci := &pb.CtrlInfo{
		CtrlID:   ctrlid,
		CtrlType: CtrlTypeScriptFile2,
		Dat:      dat,
	}

	return ci, nil
}
