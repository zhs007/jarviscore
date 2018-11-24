package jarviscore

import (
	"testing"
)

func TestIsValidNodeName(t *testing.T) {
	arrOK := []string{
		"zhs007",
		"jarviscore",
		"jarvis_dt",
		"j123_456dt",
	}

	for _, v := range arrOK {
		if !IsValidNodeName(v) {
			t.Fatalf("IsValidNodeName(%v): got false", v)
		}
	}

	arrFalse := []string{
		"007zhs",
		"",
		"a b c",
		"_aa456dt",
		"aa456dt_",
		"aa456dt ",
	}

	for _, v := range arrFalse {
		if IsValidNodeName(v) {
			t.Fatalf("IsValidNodeName(%v): got true", v)
		}
	}

	t.Log("success IsValidNodeName()")
}
