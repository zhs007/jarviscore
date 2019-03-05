package jarviscore

import (
	"testing"
)

func TestUpdateNode(t *testing.T) {
	param := &UpdateNodeParam{
		NewVersion: "0.2.6",
	}

	curscript, outstring, err := updateNode(param, "echo {{.NewVersion}}")
	if err != nil {
		t.Fatalf("TestUpdateNode err %v", err)

		return
	}

	if curscript != "echo 0.2.6" {
		t.Fatalf("TestUpdateNode curscript err (%v - %v)", curscript, "echo 0.2.6")

		return
	}

	t.Logf("TestUpdateNode result is %v", outstring)
}
