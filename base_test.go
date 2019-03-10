package jarviscore

import (
	"testing"
)

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
