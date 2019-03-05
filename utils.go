package jarviscore

import (
	"crypto/md5"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/crypto"
	pb "github.com/zhs007/jarviscore/proto"
)

// GetMD5String - md5 buf and return string
func GetMD5String(buf []byte) string {
	return fmt.Sprintf("%x", md5.Sum(buf))
}

// GetRealFilename - get filename
func GetRealFilename(fn string) string {
	arr := strings.Split(fn, "/")
	return arr[len(arr)-1]
}

// buildSignBuf - build sign buf
//		sign(msgID + msgType + destAddr + curTime + srcAddr + data)
//		for mul-language, all become string merge data
func buildSignBuf(msg *pb.JarvisMsg) ([]byte, error) {
	if msg.MsgType == pb.MSGTYPE_LOCAL_CONNECT_OTHER ||
		msg.MsgType == pb.MSGTYPE_CONNECT_NODE {

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
	} else if msg.MsgType == pb.MSGTYPE_REPLY {
		str := []byte(fmt.Sprintf("%v%v%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr, msg.ReplyType, msg.Err))

		return str, nil
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
	} else if msg.MsgType == pb.MSGTYPE_LOCAL_SENDMSG {
		cr := msg.GetMsg()
		if cr != nil {
			str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(cr)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_LOCAL_REQUEST_NODES {
		str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))

		return str, nil
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_NODES {
		str := []byte(fmt.Sprintf("%v%v%v%v%v", msg.MsgID, msg.MsgType, msg.DestAddr, msg.CurTime, msg.SrcAddr))

		return str, nil
	} else if msg.MsgType == pb.MSGTYPE_TRANSFER_FILE ||
		msg.MsgType == pb.MSGTYPE_REPLY_REQUEST_FILE {

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
	pk.FromBytes(msg.PubKey)

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

// AbsInt64 - abs int64
func AbsInt64(num int64) int64 {
	if num < 0 {
		return -num
	}

	return num
}

// IsTimeOut - is JarvisMsg timeout
func IsTimeOut(msg *pb.JarvisMsg) bool {
	ct := time.Now().Unix()
	if AbsInt64(ct-msg.CurTime) < 5*60 {
		return false
	}

	return true
}

// BuildConnNode - build jarvismsg with CONNECT_NODE
func BuildConnNode(jarvisnode JarvisNode, srcAddr string, destAddr string,
	servaddr string, ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_CONNECT_NODE,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_ConnInfo{
			ConnInfo: &pb.ConnectInfo{
				ServAddr: servaddr,
				MyInfo:   ni,
			},
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildReplyConn - build jarvismsg with REPLY_CONNECT
func BuildReplyConn(jarvisnode JarvisNode, srcAddr string, destAddr string,
	ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REPLY_CONNECT,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildLocalConnectOther - build jarvismsg with LOCAL_CONNECT_OTHER
func BuildLocalConnectOther(jarvisnode JarvisNode, srcAddr string,
	destAddr string, servaddr string, ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_LOCAL_CONNECT_OTHER,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_ConnInfo{
			ConnInfo: &pb.ConnectInfo{
				ServAddr: servaddr,
				MyInfo:   ni,
			},
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildRequestCtrl - build jarvismsg with REQUEST_CTRL
func BuildRequestCtrl(jarvisnode JarvisNode, srcAddr string,
	destAddr string, ci *pb.CtrlInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REQUEST_CTRL,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_CtrlInfo{
			CtrlInfo: ci,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildReply2 - build jarvismsg with REPLY2
func BuildReply2(jarvisnode JarvisNode, srcAddr string,
	destAddr string, rt pb.REPLYTYPE, strErr string, replyMsgID int64) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:      jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:    time.Now().Unix(),
		SrcAddr:    srcAddr,
		MyAddr:     srcAddr,
		DestAddr:   destAddr,
		MsgType:    pb.MSGTYPE_REPLY2,
		LastMsgID:  jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		ReplyType:  rt,
		Err:        strErr,
		ReplyMsgID: replyMsgID,
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildCtrlResult - build jarvismsg with REPLY_CTRL_RESULT
func BuildCtrlResult(jarvisnode JarvisNode, srcAddr string,
	destAddr string, ctrlid int64, result string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REPLY_CTRL_RESULT,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_CtrlResult{
			CtrlResult: &pb.CtrlResult{
				CtrlID:     ctrlid,
				CtrlResult: result,
			},
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildLocalSendMsg - build jarvismsg with LOCAL_SENDMSG
func BuildLocalSendMsg(jarvisnode JarvisNode, srcAddr string,
	destAddr string, sendmsg *pb.JarvisMsg) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_LOCAL_SENDMSG,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_Msg{
			Msg: sendmsg,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildLocalRequestNodes - build jarvismsg with LOCAL_REQUEST_NODES
func BuildLocalRequestNodes(jarvisnode JarvisNode, srcAddr string,
	destAddr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		MsgType:   pb.MSGTYPE_LOCAL_REQUEST_NODES,
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildRequestNodes - build jarvismsg with REQUEST_NODES
func BuildRequestNodes(jarvisnode JarvisNode, srcAddr string,
	destAddr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REQUEST_NODES,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildNodeInfo - build jarvismsg with NODE_INFO
func BuildNodeInfo(jarvisnode JarvisNode, srcAddr string, destAddr string,
	ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_NODE_INFO,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildTransferFile - build jarvismsg with TRANSFER_FILE
func BuildTransferFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	fd *pb.FileData) (*pb.JarvisMsg, error) {

	fd.Md5String = GetMD5String(fd.File)

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_TRANSFER_FILE,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildRequestFile - build jarvismsg with REQUEST_FILE
func BuildRequestFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	rf *pb.RequestFile) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REQUEST_FILE,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_RequestFile{
			RequestFile: rf,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildReplyRequestFile - build jarvismsg with REPLY_REQUEST_FILE
func BuildReplyRequestFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	fd *pb.FileData) (*pb.JarvisMsg, error) {

	fd.Md5String = GetMD5String(fd.File)

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REPLY_REQUEST_FILE,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// IsValidNodeName - check node name
func IsValidNodeName(nodename string) bool {
	if len(nodename) == 0 {
		return false
	}

	for i, v := range nodename {
		if i == 0 && ((v >= '0' && v <= '9') || v == '_') {
			return false
		}

		if !((v >= '0' && v <= '9') || (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || v == '_') {
			return false
		}
	}

	if nodename[len(nodename)-1] == '_' {
		return false
	}

	return true
}

// StoreLocalFile - store filedata to local file systems
func StoreLocalFile(file *pb.FileData) error {
	f, err := os.Create(file.Filename)
	if err != nil {
		jarvisbase.Warn("StoreLocalFile:os.Create err", zap.Error(err))

		return err
	}

	defer f.Close()

	f.Write(file.File)
	f.Close()

	return nil
}

// GetNodeBaseInfo - get nodebaseinfo from nodeinfo
func GetNodeBaseInfo(node *coredbpb.NodeInfo) *pb.NodeBaseInfo {
	return &pb.NodeBaseInfo{
		ServAddr:        node.ServAddr,
		Addr:            node.Addr,
		Name:            node.Name,
		NodeTypeVersion: node.NodeTypeVersion,
		NodeType:        node.NodeType,
		CoreVersion:     node.CoreVersion,
	}
}

// BuildReplyTransferFile - build jarvismsg with REPLY_TRANSFER_FILE
func BuildReplyTransferFile(jarvisnode JarvisNode, srcAddr string, destAddr string,
	md5str string, replyMsgID int64) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:      jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:    time.Now().Unix(),
		SrcAddr:    srcAddr,
		MyAddr:     srcAddr,
		DestAddr:   destAddr,
		MsgType:    pb.MSGTYPE_REPLY_TRANSFER_FILE,
		LastMsgID:  jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		ReplyMsgID: replyMsgID,
		Data: &pb.JarvisMsg_ReplyTransferFile{
			ReplyTransferFile: &pb.ReplyTransferFile{
				Md5String: md5str,
			},
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildUpdateNode - build jarvismsg with UPDATENODE
func BuildUpdateNode(jarvisnode JarvisNode, srcAddr string, destAddr string,
	nodetype string, nodetypever string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     jarvisnode.GetCoreDB().GetNewSendMsgID(destAddr),
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_UPDATENODE,
		LastMsgID: jarvisnode.GetCoreDB().GetCurRecvMsgID(destAddr),
		Data: &pb.JarvisMsg_UpdateNode{
			UpdateNode: &pb.UpdateNode{
				NodeType:        nodetype,
				NodeTypeVersion: nodetypever,
			},
		},
	}

	err := SignJarvisMsg(jarvisnode.GetCoreDB().GetPrivateKey(), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
