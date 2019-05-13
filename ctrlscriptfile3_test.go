package jarviscore

import (
	"io/ioutil"
	"testing"

	pb "github.com/zhs007/jarviscore/proto"
)

func TestCtrlScriptFile3(t *testing.T) {
	ctrl := &CtrlScriptFile3{}

	dat, err := ioutil.ReadFile("./test/test.sh")
	if err != nil {
		t.Fatalf("TestCtrlScriptFile3 load script file %v", err)
	}

	sf := &pb.FileData{
		Filename: "test.sh",
		File:     dat,
	}

	ci, err := BuildCtrlInfoForScriptFile3(sf, nil, "")
	if err != nil {
		t.Fatalf("TestCtrlScriptFile3 BuildCtrlInfoForScriptFile3 %v", err)
	}

	_, ret, err := ctrl.runScript(nil, nil, ci)
	if err != nil {
		t.Fatalf("TestCtrlScriptFile3 Error Run %v", err)
	}

	t.Logf("TestCtrlScriptFile3 result is %v", string(ret))
}
