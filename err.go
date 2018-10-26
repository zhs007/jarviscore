package jarviscore

import "errors"

var (
	// ErrLoadFileReadSize - loadfile invalid file read size
	ErrLoadFileReadSize = errors.New("loadfile invalid file read size")
	// ErrNotConnectNode - not connect node
	ErrNotConnectNode = errors.New("not connect node")
	// ErrNoPrivateKey - no private key
	ErrNoPrivateKey = errors.New("no private key")
	// ErrNoCtrlCmd - no ctrl cmd
	ErrNoCtrlCmd = errors.New("no ctrl cmd")
	// ErrCoreDBNoAddr - coredb no addr
	ErrCoreDBNoAddr = errors.New("coredb no addr")
	// ErrSign - sign err
	ErrSign = errors.New("sign err")
	// ErrInvalidAddr - invalid addr
	ErrInvalidAddr = errors.New("invalid addr")
	// ErrInvalidPublishKey - invalid publish key
	ErrInvalidPublishKey = errors.New("invalid publish key")
	// ErrExistCtrlID - exist ctrlid
	ErrExistCtrlID = errors.New("exist ctrlid")
	// ErrAlreadyJoin - already join
	ErrAlreadyJoin = errors.New("already join")
	// ErrNotConnectMe - not connect me
	ErrNotConnectMe = errors.New("not connect me")
)
