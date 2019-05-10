package jarviscore

import (
	"context"
	"errors"

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
func (ctrl *CtrlScriptFile) runScript(logpath string, ci *pb.CtrlInfo) ([]byte, error) {
	var csd pb.CtrlScriptData
	err := ptypes.UnmarshalAny(ci.Dat, &csd)
	if err != nil {
		return nil, err
	}

	outstr, errstr, err := RunCommand(logpath, string(csd.File))
	if err != nil {
		return nil, err
	}

	if errstr == "" {
		return []byte(outstr), nil
	}

	return []byte(outstr), errors.New(errstr)
}

// Run -
func (ctrl *CtrlScriptFile) Run(ctx context.Context, jarvisnode JarvisNode, srcAddr string, msgid int64, ci *pb.CtrlInfo) []*pb.JarvisMsg {

	var msgs []*pb.JarvisMsg

	out, err := ctrl.runScript(jarvisnode.GetConfig().Log.LogPath, ci)
	if err != nil {
		if out == nil {
			return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, err.Error(), msgs)
		}

		return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, AppendString(string(out), err.Error()), msgs)
	}

	return BuildCtrlResultForCtrl(jarvisnode, srcAddr, msgid, string(out), msgs)
}

// BuildCtrlInfoForScriptFile - build ctrlinfo for scriptfile
// Deprecated: you can use BuildCtrlInfoForScriptFile3
func BuildCtrlInfoForScriptFile(filename string, filedata []byte,
	destpath string, scriptName string) (*pb.CtrlInfo, error) {

	csd := &pb.CtrlScriptData{
		File:       filedata,
		DestPath:   destpath,
		Filename:   filename,
		ScriptName: scriptName,
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
