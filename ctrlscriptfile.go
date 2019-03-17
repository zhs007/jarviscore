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

// Run -
func (ctrl *CtrlScriptFile) Run(ci *pb.CtrlInfo) ([]byte, error) {
	var csd pb.CtrlScriptData
	err := ptypes.UnmarshalAny(ci.Dat, &csd)
	if err != nil {
		return nil, err
	}

	out, err := exec.Command("sh", "-c", string(csd.File)).CombinedOutput()
	if err != nil {
		return out, err
	}

	return out, nil
}

// BuildCtrlInfoForScriptFile - build ctrlinfo for scriptfile
func BuildCtrlInfoForScriptFile(ctrlid int64, filename string, filedata []byte,
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
		CtrlID:   ctrlid,
		CtrlType: CtrlTypeScriptFile,
		Dat:      dat,
	}

	return ci, nil
}
