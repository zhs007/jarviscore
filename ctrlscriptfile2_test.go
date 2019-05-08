package jarviscore

import (
	"io/ioutil"
	"testing"

	pb "github.com/zhs007/jarviscore/proto"
)

func TestCtrlScriptFile2(t *testing.T) {
	ctrl := &CtrlScriptFile2{}

	dat, err := ioutil.ReadFile("./test/test.sh")
	if err != nil {
		t.Fatalf("TestCtrlScriptFile2 load script file %v", err)
	}

	sf := &pb.FileData{
		Filename: "test.sh",
		File:     dat,
	}

	ci, err := BuildCtrlInfoForScriptFile2(sf, nil)
	if err != nil {
		t.Fatalf("TestCtrlScriptFile2 BuildCtrlInfoForScriptFile2 %v", err)
	}

	ret, err := ctrl.runScript("./", ci)
	if err != nil {
		t.Fatalf("TestCtrlScriptFile2 Error Run %v", err)
	}

	t.Logf("TestCtrlScriptFile2 result is %v", string(ret))
}
