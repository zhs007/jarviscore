package jarviscore

import (
	"time"

	"github.com/zhs007/jarviscore/basedef"
)

// failservaddrinfo - this servaddr is fail, we won't try to connect it for a while
type failservaddrinfo struct {
	timestampConnFail int64
	numsConnFail      int
}

// failservaddr -
type failservaddr struct {
	mapServAddr map[string]*failservaddrinfo
}

func newFailServAddr() *failservaddr {
	return &failservaddr{
		mapServAddr: make(map[string]*failservaddrinfo),
	}
}

func (fsa *failservaddr) onConnFail(servaddr string) {
	sa, ok := fsa.mapServAddr[servaddr]
	if ok {
		sa.numsConnFail++

		return
	}

	fsa.mapServAddr[servaddr] = &failservaddrinfo{
		timestampConnFail: time.Now().Unix(),
		numsConnFail:      1,
	}
}

func (fsa *failservaddr) isFailServAddr(servaddr string) bool {
	sa, ok := fsa.mapServAddr[servaddr]
	if ok {
		if sa.numsConnFail > basedef.NumsReconnectFailedServAddr {
			timeUnix := time.Now().Unix()
			if timeUnix <= basedef.TimeFailedServAddr+sa.timestampConnFail {
				return true
			}

			delete(fsa.mapServAddr, servaddr)
		}
	}

	return false
}
