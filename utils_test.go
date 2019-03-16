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

func TestIsMyServAddr(t *testing.T) {
	type data struct {
		destaddr string
		srcaddr  string
		ret      bool
	}

	lst := []data{
		data{"192.168.0.1:7788", "192.168.0.1:7788", true},
		data{"192.168.0.1:7788", "127.0.0.1:7788", false},
		data{"127.0.0.1:7788", "192.168.0.1:7788", true},
		data{"192.168.0.1:7788", "127.0.0.1:7789", false},
		data{"127.0.0.1:7789", "192.168.0.1:7788", false},
	}

	for i := 0; i < len(lst); i++ {
		cr := IsMyServAddr(lst[i].destaddr, lst[i].srcaddr)
		if cr != lst[i].ret {
			t.Fatalf("TestIsMyServAddr fail %v %v", lst[i].destaddr, lst[i].srcaddr)
		}
	}

	t.Logf("TestIsMyServAddr OK")
}

func TestIsValidServAddr(t *testing.T) {
	type data struct {
		servaddr string
		ret      bool
	}

	lst := []data{
		data{"192.168.0.1:7788", true},
		data{"192.168.0.1:", false},
		data{":7788", false},
		data{"a.b.c:7788", true},
	}

	for i := 0; i < len(lst); i++ {
		cr := IsValidServAddr(lst[i].servaddr)
		if cr != lst[i].ret {
			t.Fatalf("TestIsValidServAddr fail %v", lst[i].servaddr)
		}
	}

	t.Logf("TestIsValidServAddr OK")
}

func TestIsLocalHostAddr(t *testing.T) {
	type data struct {
		servaddr string
		ret      bool
	}

	lst := []data{
		data{"192.168.0.1:7788", false},
		data{"192.168.0.1:", false},
		data{":7788", false},
		data{"localhost:7788", true},
		data{"localhost:", true},
		data{"127.0.0.1:7788", true},
		data{"127.0.0.1:", true},
	}

	for i := 0; i < len(lst); i++ {
		cr := IsLocalHostAddr(lst[i].servaddr)
		if cr != lst[i].ret {
			t.Fatalf("TestIsLocalHostAddr fail %v", lst[i].servaddr)
		}
	}

	t.Logf("TestIsLocalHostAddr OK")
}
