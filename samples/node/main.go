package main

import (
	"github.com/zhs007/jarviscore"
)

func main() {
	node := jarviscore.NewNode()
	node.Start("127.0.0.1:7788", "node001", "")
}
