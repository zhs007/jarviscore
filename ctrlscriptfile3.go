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
func (ctrl *CtrlScriptFile3) runScript(ci *pb.CtrlInfo) (*pb.CtrlScript3Data, []byte, error) {
	var csd3 pb.CtrlScript3Data
	err := ptypes.UnmarshalAny(ci.Dat, &csd3)
	if err != nil {
		return nil, nil, err
	}

	out, err := exec.Command("sh", "-c", string(csd3.ScriptFile.File)).CombinedOutput()

	return &csd3, out, err
}

// Run -
func (ctrl *CtrlScriptFile3) Run(jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {

	var msgs []*pb.JarvisMsg

	csd3, out, err := ctrl.runScript(ci)
	if err != nil {
		return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, AppendString(string(out), err.Error()), msgs)
	}

	msgs = BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, string(out), msgs)

	for i := 0; i < len(csd3.EndFiles); i++ {
		err := ProcFileData(csd3.EndFiles[i], func(fd *pb.FileData, isend bool) error {
			msgs = BuildReplyRequestFileForCtrl(jarvisnode, srcAddr, msgid, fd, msgs)
			// sendmsg, err := BuildReplyRequestFile(jarvisnode, jarvisnode.GetMyInfo().Addr, srcAddr, fd, msgid)
			// if err != nil {
			// 	jarvisbase.Warn("CtrlScriptFile3.Run", zap.Error(err))
			// }

			// msgs = append(msgs, sendmsg)

			return nil
		})
		if err != nil {
			msgs = BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, err.Error(), msgs)
		}
	}

	return msgs
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
