package jarviscore

import (
	"io/ioutil"
	"testing"
)

func TestCtrlScriptFile(t *testing.T) {
	ctrl := &CtrlScriptFile{}

	dat, err := ioutil.ReadFile("./test/test.sh")
	if err != nil {
		t.Fatalf("TestCtrlScriptFile load script file %v", err)
	}

	ci, err := BuildCtrlInfoForScriptFile("test.sh", dat, "", "")
	if err != nil {
		t.Fatalf("TestCtrlScriptFile BuildCtrlInfoForScriptFile %v", err)
	}

	ret, err := ctrl.runScript(nil, nil, ci)
	if err != nil {
		t.Fatalf("TestCtrlScriptFile Run %v", err)
	}

	t.Logf("TestCtrlScriptFile result is %v", string(ret))
}
