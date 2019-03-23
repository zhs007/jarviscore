package jarviscore

import "errors"

var (
	// ErrLoadFileReadSize - loadfile invalid file read size
	ErrLoadFileReadSize = errors.New("loadfile invalid file read size")
	// ErrNotConnectNode - not connect node
	ErrNotConnectNode = errors.New("not connect node")
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
	// ErrInvalidNodeName - invalid nodename
	ErrInvalidNodeName = errors.New("invalid nodename")
	// ErrUnknowNode - unknow node
	ErrUnknowNode = errors.New("unknow node")
	// ErrInvalidMsgID - invalid msgid
	ErrInvalidMsgID = errors.New("invalid msgid")
	// ErrDuplicateMsgID - duplicate msgid
	ErrDuplicateMsgID = errors.New("duplicate msgid")
	// ErrInvalidRequestData4Node - invalid requestData4Node
	ErrInvalidRequestData4Node = errors.New("invalid requestData4Node")
	// ErrInvalidRequestNodeData - invalid requestNodeData
	ErrInvalidRequestNodeData = errors.New("invalid requestNodeData")
	// ErrFuncOnSendMsgResultLength - FuncOnSendMsgResult length err
	ErrFuncOnSendMsgResultLength = errors.New("FuncOnSendMsgResult length err")
	// ErrAutoUpdateClosed - auto update closed
	ErrAutoUpdateClosed = errors.New("auto update closed")
	// ErrServAddrConnFail - the servaddr connect fail
	ErrServAddrConnFail = errors.New("the servaddr connect fail")
	// ErrAssertGetNode - assert(GetNode() is not nil)
	ErrAssertGetNode = errors.New("assert(GetNode() is not nil)")
	// ErrDeprecatedNode - deprecated node
	ErrDeprecatedNode = errors.New("deprecated node")
	// ErrNoFileData - no filedata
	ErrNoFileData = errors.New("no filedata")
	// ErrFileDataNoMD5String - filedata no md5tring
	ErrFileDataNoMD5String = errors.New("filedata no md5tring")
	// ErrInvalidFileDataMD5String - invalid filedata md5tring
	ErrInvalidFileDataMD5String = errors.New("invalid filedata md5tring")
	// ErrNotConnectedNode - not connected node
	ErrNotConnectedNode = errors.New("not connected node")
	// ErrInvalidReadFileLength - invalid readfile length
	ErrInvalidReadFileLength = errors.New("invalid readfile length")
	// ErrInvalidSeekFileOffset - invalid seekfile offset
	ErrInvalidSeekFileOffset = errors.New("invalid seekfile offset")
	// ErrNoConnOrInvalidConn - no connection or invalid connection
	ErrNoConnOrInvalidConn = errors.New("no connection or invalid connection")
	// ErrNoCtrlInfo - no ctrlinfo
	ErrNoCtrlInfo = errors.New("no ctrlinfo")
	// ErrCannotFindNodeWithAddr - can not find node with addr
	ErrCannotFindNodeWithAddr = errors.New("can not find node with addr")
	// ErrNoFuncOnFileData - no FuncOnFileData
	ErrNoFuncOnFileData = errors.New("no FuncOnFileData")
	// ErrNoProcMsgResultData - no ProcMsgResultData
	ErrNoProcMsgResultData = errors.New("no ProcMsgResultData")
	// ErrInvalidProcMsgResultData - invalid ProcMsgResultData
	ErrInvalidProcMsgResultData = errors.New("invalid ProcMsgResultData")
	// ErrDuplicateProcMsgResultData - duplicate ProcMsgResultData
	ErrDuplicateProcMsgResultData = errors.New("duplicate ProcMsgResultData")
	// ErrNoCtrl - no ctrl
	ErrNoCtrl = errors.New("no ctrl")
	// ErrUnknownCtrlError - unknown ctrl error
	ErrUnknownCtrlError = errors.New("unknown ctrl error")
	// ErrProcMsgStreamNil - ProcMsgStream return nil
	ErrProcMsgStreamNil = errors.New("ProcMsgStream return nil")
)
