package jarviscore

import (
	"fmt"
	"math/big"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	jarvisbase "github.com/zhs007/jarviscore/base"
	jarviscrypto "github.com/zhs007/jarviscore/crypto"
	pb "github.com/zhs007/jarviscore/proto"
)

// buildSignBuf - build sign buf
//		sign(msgID + msgType + destAddr + curTime + srcAddr + data)
//		for mul-language, all become string merge data
func buildSignBuf(msg *pb.JarvisMsg) ([]byte, error) {
	if msg.MsgType == pb.MSGTYPE_CONNECT_NODE {

		ci := msg.GetConnInfo()
		if ci != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(ci)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_NODE_INFO ||
		msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {

		ni := msg.GetNodeInfo()
		if ni != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(ni)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_CTRL {

		ci := msg.GetCtrlInfo()
		if ci != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(ci)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REPLY_CTRL_RESULT {
		cr := msg.GetCtrlResult()
		if cr != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(cr)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_NODES ||
		msg.MsgType == pb.MSGTYPE_CLEAR_LOGS ||
		msg.MsgType == pb.MSGTYPE_REQUEST_LOGLIST {

		str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))

		return str, nil
	} else if msg.MsgType == pb.MSGTYPE_TRANSFER_FILE ||
		msg.MsgType == pb.MSGTYPE_REPLY_REQUEST_FILE ||
		msg.MsgType == pb.MSGTYPE_TRANSFER_FILE2 {

		f := msg.GetFile()
		if f != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(f)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_FILE {
		rf := msg.GetRequestFile()
		if rf != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(rf)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REPLY_TRANSFER_FILE {
		rtf := msg.GetReplyTransferFile()
		if rtf != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr,
				msg.ReplyMsgID))
			buf, err := proto.Marshal(rtf)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REPLY2 {
		str := []byte(fmt.Sprintf("%v%v%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr,
			msg.ReplyType, msg.Err, msg.ReplyMsgID))

		return str, nil
	} else if msg.MsgType == pb.MSGTYPE_UPDATENODE {
		rf := msg.GetUpdateNode()
		if rf != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(rf)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_MULTI_MSG {
		mmd := msg.GetMultiMsgData()
		if mmd != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr,
				msg.ReplyMsgID))
			buf, err := proto.Marshal(mmd)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_MSG_STATE {
		rms := msg.GetRequestMsgState()
		if rms != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(rms)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REPLY_MSG_STATE {
		rms := msg.GetReplyMsgState()
		if rms != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr,
				msg.ReplyMsgID))
			buf, err := proto.Marshal(rms)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_NODES2 {
		rn2 := msg.GetRequestNodes2()
		if rn2 != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr,
				msg.ReplyMsgID))
			buf, err := proto.Marshal(rn2)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REPLY_MYNODESVERSION {
		rmnv := msg.GetReplyMyNodesVersion()
		if rmnv != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr,
				msg.ReplyMsgID))
			buf, err := proto.Marshal(rmnv)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	}

	return nil, ErrInvalidMsgType
}

// SignJarvisMsg - sign JarvisMsg
func SignJarvisMsg(privkey *jarviscrypto.PrivateKey, msg *pb.JarvisMsg) error {
	buf, err := buildSignBuf(msg)
	if err != nil {
		return err
	}

	r, s, err := privkey.Sign(buf)
	if err != nil {
		return ErrSign
	}

	msg.PubKey = privkey.ToPublicBytes()
	msg.SignR = r.Bytes()
	msg.SignS = s.Bytes()

	return nil
}

// VerifyJarvisMsg - Verify JarvisMsg
func VerifyJarvisMsg(msg *pb.JarvisMsg) error {
	pk := jarviscrypto.PublicKey{}
	err := pk.FromBytes(msg.PubKey)
	if err != nil {
		return err
	}

	if pk.ToAddress() != msg.SrcAddr {
		return ErrPublicKeyAddr
	}

	buf, err := buildSignBuf(msg)
	if err != nil {
		return err
	}

	s := new(big.Int)
	r := new(big.Int)

	s.SetBytes(msg.SignS)
	r.SetBytes(msg.SignR)

	if !pk.Verify(buf, r, s) {
		return ErrPublicKeyVerify
	}

	return nil
}

// BuildConnNode - build jarvismsg with CONNECT_NODE
func BuildConnNode(jarvisnode JarvisNode, srcAddr string, destAddr string,
	servaddr string, ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_CONNECT_NODE,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_ConnInfo{
			ConnInfo: &pb.ConnectInfo{
				ServAddr: servaddr,
				MyInfo:   ni,
			},
		},
	}

	return msg, nil
}

// BuildReplyConn - build jarvismsg with REPLY_CONNECT
func BuildReplyConn(jarvisnode JarvisNode, srcAddr string, destAddr string,
	ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_CONNECT,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	return msg, nil
}

// BuildRequestCtrl - build jarvismsg with REQUEST_CTRL
func BuildRequestCtrl(jarvisnode JarvisNode, srcAddr string,
	destAddr string, ci *pb.CtrlInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_CTRL,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_CtrlInfo{
			CtrlInfo: ci,
		},
	}

	return msg, nil
}

// BuildReply2 - build jarvismsg with REPLY2
func BuildReply2(jarvisnode JarvisNode, srcAddr string,
	destAddr string, rt pb.REPLYTYPE, strErr string, replyMsgID int64) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY2,
		// LastMsgID:  jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		ReplyType:  rt,
		Err:        strErr,
		ReplyMsgID: replyMsgID,
	}

	return msg, nil
}

// BuildCtrlResult - build jarvismsg with REPLY_CTRL_RESULT
func BuildCtrlResult(jarvisnode JarvisNode, srcAddr string,
	destAddr string, replyMsgID int64, result string,
	errInfo string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_CTRL_RESULT,
		// LastMsgID:  jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		ReplyMsgID: replyMsgID,
		Data: &pb.JarvisMsg_CtrlResult{
			CtrlResult: &pb.CtrlResult{
				CtrlResult: result,
				ErrInfo:    errInfo,
			},
		},
	}

	return msg, nil
}

// BuildRequestNodes - build jarvismsg with REQUEST_NODES
func BuildRequestNodes(jarvisnode JarvisNode, srcAddr string,
	destAddr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_NODES,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
	}

	return msg, nil
}

// BuildRequestNodes2 - build jarvismsg with REQUEST_NODES2
func BuildRequestNodes2(jarvisnode JarvisNode, srcAddr string,
	destAddr string, isNeedLocalHost bool, myNodesVersion string, nodesVersion string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_NODES2,
		Data: &pb.JarvisMsg_RequestNodes2{
			RequestNodes2: &pb.RequestNodes2{
				IsNeedLocalHost: isNeedLocalHost,
				MyNodesVersion:  myNodesVersion,
				NodesVersion:    nodesVersion,
			},
		},
	}

	return msg, nil
}

// BuildNodeInfo - build jarvismsg with NODE_INFO
func BuildNodeInfo(jarvisnode JarvisNode, srcAddr string, destAddr string,
	ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_NODE_INFO,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	return msg, nil
}

// BuildTransferFile - build jarvismsg with TRANSFER_FILE
func BuildTransferFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	fd *pb.FileData) (*pb.JarvisMsg, error) {

	fd.Md5String = GetMD5String(fd.File)

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_TRANSFER_FILE,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	return msg, nil
}

// BuildTransferFile2 - build jarvismsg with TRANSFER_FILE2
func BuildTransferFile2(jarvisnode JarvisNode, srcAddr string, destAddr string,
	fd *pb.FileData) (*pb.JarvisMsg, error) {

	// fd.Md5String = GetMD5String(fd.File)

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_TRANSFER_FILE2,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	return msg, nil
}

// BuildRequestFile - build jarvismsg with REQUEST_FILE
func BuildRequestFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	rf *pb.RequestFile) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_FILE,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_RequestFile{
			RequestFile: rf,
		},
	}

	return msg, nil
}

// BuildReplyRequestFile - build jarvismsg with REPLY_REQUEST_FILE
func BuildReplyRequestFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	fd *pb.FileData, replyMsgID int64) (*pb.JarvisMsg, error) {

	// fd.Md5String = GetMD5String(fd.File)

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_REQUEST_FILE,
		// LastMsgID:  jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		ReplyMsgID: replyMsgID,
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	return msg, nil
}

// BuildReplyTransferFile - build jarvismsg with REPLY_TRANSFER_FILE
func BuildReplyTransferFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	md5str string, replyMsgID int64) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_TRANSFER_FILE,
		// LastMsgID:  jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		ReplyMsgID: replyMsgID,
		Data: &pb.JarvisMsg_ReplyTransferFile{
			ReplyTransferFile: &pb.ReplyTransferFile{
				Md5String: md5str,
			},
		},
	}

	return msg, nil
}

// BuildUpdateNode - build jarvismsg with UPDATENODE
func BuildUpdateNode(jarvisnode JarvisNode, srcAddr string, destAddr string,
	nodetype string, nodetypever string, isOnlyRestart bool) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_UPDATENODE,
		// LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_UpdateNode{
			UpdateNode: &pb.UpdateNode{
				NodeType:        nodetype,
				NodeTypeVersion: nodetypever,
				IsOnlyRestart:   isOnlyRestart,
			},
		},
	}

	return msg, nil
}

// BuildOutputMsg - build a output jarvismsg
func BuildOutputMsg(src *pb.JarvisMsg) (*pb.JarvisMsg, error) {
	if src.MsgType == pb.MSGTYPE_REPLY_REQUEST_FILE || src.MsgType == pb.MSGTYPE_TRANSFER_FILE {
		dest := &pb.JarvisMsg{}

		err := DeepCopy(src, dest)
		if err != nil {
			return nil, err
		}

		dest.GetFile().File = nil

		return dest, nil
	}

	return src, nil
}

// JSONMsg2Zap - I use this interface to output jarvismsg to the zap log.
//		This interface will hide the long data in jarvismsg.
func JSONMsg2Zap(key string, src *pb.JarvisMsg) zap.Field {
	msg, err := BuildOutputMsg(src)
	if err != nil {
		jarvisbase.Warn("JSONMsg2Zap:BuildOutputMsg", zap.Error(err))

		return jarvisbase.JSON(key, src)
	}

	return jarvisbase.JSON(key, msg)
}

// PushReply22Msgs - push Reply2 to msgs
func PushReply22Msgs(msgs []*pb.JarvisMsg, jarvisnode JarvisNode, srcAddr string, msgid int64,
	replytype pb.REPLYTYPE, info string) []*pb.JarvisMsg {

	msg, err1 := BuildReply2(jarvisnode,
		jarvisnode.GetMyInfo().Addr,
		srcAddr,
		replytype,
		info,
		msgid)

	if err1 != nil {
		jarvisbase.Warn("PushReply22Msgs", zap.Error(err1))

		return msgs
	}

	return append(msgs, msg)
}

// NewErrorMsg - new a error JarvisMsg
func NewErrorMsg(jarvisnode JarvisNode, nodeAddr string, strErr string, replyMsgID int64) *pb.JarvisMsg {
	return &pb.JarvisMsg{
		CurTime:    time.Now().Unix(),
		SrcAddr:    jarvisnode.GetMyInfo().Addr,
		MyAddr:     jarvisnode.GetMyInfo().Addr,
		DestAddr:   nodeAddr,
		MsgType:    pb.MSGTYPE_REPLY2,
		ReplyType:  pb.REPLYTYPE_ERROR,
		Err:        strErr,
		ReplyMsgID: replyMsgID,
	}
}

// NewCtrlResult - new jarvismsg with REPLY_CTRL_RESULT
func NewCtrlResult(jarvisnode JarvisNode, nodeAddr string, msgid int64, dat proto.Message) (*pb.JarvisMsg, error) {

	anydat, err := ptypes.MarshalAny(dat)
	if err != nil {
		return nil, err
	}

	msg := &pb.JarvisMsg{
		CurTime:    time.Now().Unix(),
		SrcAddr:    jarvisnode.GetMyInfo().Addr,
		MyAddr:     jarvisnode.GetMyInfo().Addr,
		DestAddr:   nodeAddr,
		MsgType:    pb.MSGTYPE_REPLY_CTRL_RESULT,
		ReplyMsgID: msgid,
		Data: &pb.JarvisMsg_CtrlResult{
			CtrlResult: &pb.CtrlResult{
				Dat: anydat,
			},
		},
	}

	return msg, nil
}

// BuildMultiMsgData - build jarvismsg with MULTI_MSG
func BuildMultiMsgData(jarvisnode JarvisNode, destAddr string, replyMsgID int64,
	multimsg *pb.MultiMsgData) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:    time.Now().Unix(),
		SrcAddr:    jarvisnode.GetMyInfo().Addr,
		MyAddr:     jarvisnode.GetMyInfo().Addr,
		DestAddr:   destAddr,
		MsgType:    pb.MSGTYPE_MULTI_MSG,
		ReplyMsgID: replyMsgID,
		Data: &pb.JarvisMsg_MultiMsgData{
			MultiMsgData: multimsg,
		},
	}

	return msg, nil
}

// BuildRequestMsgState - build jarvismsg with REQUEST_MSG_STATE
func BuildRequestMsgState(jarvisnode JarvisNode, destAddr string,
	msgid int64) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  jarvisnode.GetMyInfo().Addr,
		MyAddr:   jarvisnode.GetMyInfo().Addr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_MSG_STATE,
		Data: &pb.JarvisMsg_RequestMsgState{
			RequestMsgState: &pb.RequestMsgState{
				MsgID: msgid,
			},
		},
	}

	return msg, nil
}

// BuildReplyMsgState - build jarvismsg with REPLY_MSG_STATE
func BuildReplyMsgState(jarvisnode JarvisNode, destAddr string, replyMsgID int64,
	msgid int64, state int) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:    time.Now().Unix(),
		SrcAddr:    jarvisnode.GetMyInfo().Addr,
		MyAddr:     jarvisnode.GetMyInfo().Addr,
		DestAddr:   destAddr,
		MsgType:    pb.MSGTYPE_REPLY_MSG_STATE,
		ReplyMsgID: replyMsgID,
		Data: &pb.JarvisMsg_ReplyMsgState{
			ReplyMsgState: &pb.ReplyMsgState{
				MsgID: msgid,
				State: int32(state),
			},
		},
	}

	return msg, nil
}

// IsSyncMsg - Is it a sync message?
func IsSyncMsg(msg *pb.JarvisMsg) bool {
	return msg.MsgType == pb.MSGTYPE_CONNECT_NODE || msg.MsgType == pb.MSGTYPE_REQUEST_MSG_STATE
}

// BuildClearLogs - build jarvismsg with CLEAR_LOGS
func BuildClearLogs(jarvisnode JarvisNode, srcAddr string,
	destAddr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_CLEAR_LOGS,
	}

	return msg, nil
}

// BuildReplyMyNodesVersion - build jarvismsg with REPLY_MYNODESVERSION
func BuildReplyMyNodesVersion(jarvisnode JarvisNode, srcAddr string,
	destAddr string, myNodesVersion string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_MYNODESVERSION,
		Data: &pb.JarvisMsg_ReplyMyNodesVersion{
			ReplyMyNodesVersion: &pb.ReplyMyNodesVersion{
				MyNodesVersion: myNodesVersion,
			},
		},
	}

	return msg, nil
}
