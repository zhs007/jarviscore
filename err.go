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
	// ErrGRPCPeerFromContext - grpc.peer.FromContext err
	ErrGRPCPeerFromContext = errors.New("grpc.peer.FromContext err")
	// ErrGRPCPeerAddr - grpc.peer.Addr err
	ErrGRPCPeerAddr = errors.New("grpc.peer.Addr err")
	// ErrPublicKeyAddr - public key and address do not match
	ErrPublicKeyAddr = errors.New("public key and address do not match")
	// ErrPublicKeyVerify - public key verify err
	ErrPublicKeyVerify = errors.New("public key verify err")
	// ErrJarvisMsgTimeOut - JarvisMsg timeout
	ErrJarvisMsgTimeOut = errors.New("JarvisMsg timeout")
	// ErrCoreDBHasNotNode - coredb has not node
	ErrCoreDBHasNotNode = errors.New("coredb has not node")
	// ErrStreamNil - stream nil
	ErrStreamNil = errors.New("stream nil")
	// ErrInvalidMsgType - invalid msgtype
	ErrInvalidMsgType = errors.New("invalid msgtype")
	// ErrServAddrIsMe - servaddr is me
	ErrServAddrIsMe = errors.New("servaddr is me")
	// ErrInvalidEvent - invalid event
	ErrInvalidEvent = errors.New("invalid event")
	// ErrInvalidServAddr - invalid servaddr
	ErrInvalidServAddr = errors.New("invalid servaddr")
)
