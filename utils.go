package jarviscore

import (
	"fmt"
	"math/big"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/zhs007/jarviscore/crypto"
	pb "github.com/zhs007/jarviscore/proto"
)

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
	}

	// jarvisbase.Debug("buildSignBuf", zap.Error(ErrInvalidMsgType))

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
func BuildConnNode(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string, destAddr string,
	servaddr string, ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_CONNECT_NODE,
		Data: &pb.JarvisMsg_ConnInfo{
			ConnInfo: &pb.ConnectInfo{
				ServAddr: servaddr,
				MyInfo:   ni,
			},
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildReplyConn - build jarvismsg with REPLY_CONNECT
func BuildReplyConn(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string, destAddr string,
	ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_CONNECT,
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildLocalConnectOther - build jarvismsg with LOCAL_CONNECT_OTHER
func BuildLocalConnectOther(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string, servaddr string, ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_LOCAL_CONNECT_OTHER,
		Data: &pb.JarvisMsg_ConnInfo{
			ConnInfo: &pb.ConnectInfo{
				ServAddr: servaddr,
				MyInfo:   ni,
			},
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildRequestCtrl - build jarvismsg with REQUEST_CTRL
func BuildRequestCtrl(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string, ci *pb.CtrlInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_CTRL,
		Data: &pb.JarvisMsg_CtrlInfo{
			CtrlInfo: ci,
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildReply - build jarvismsg with REPLY
func BuildReply(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string, rt pb.REPLYTYPE, strErr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:     msgid,
		CurTime:   time.Now().Unix(),
		SrcAddr:   srcAddr,
		MyAddr:    srcAddr,
		DestAddr:  destAddr,
		MsgType:   pb.MSGTYPE_REPLY,
		ReplyType: rt,
		Err:       strErr,
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildCtrlResult - build jarvismsg with REPLY_CTRL_RESULT
func BuildCtrlResult(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string, ctrlid int64, result string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_CTRL_RESULT,
		Data: &pb.JarvisMsg_CtrlResult{
			CtrlResult: &pb.CtrlResult{
				CtrlID:     ctrlid,
				CtrlResult: result,
			},
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildLocalSendMsg - build jarvismsg with LOCAL_SENDMSG
func BuildLocalSendMsg(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string, sendmsg *pb.JarvisMsg) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_LOCAL_SENDMSG,
		Data: &pb.JarvisMsg_Msg{
			Msg: sendmsg,
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildLocalRequestNodes - build jarvismsg with LOCAL_REQUEST_NODES
func BuildLocalRequestNodes(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_LOCAL_REQUEST_NODES,
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildRequestNodes - build jarvismsg with REQUEST_NODES
func BuildRequestNodes(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string,
	destAddr string) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_NODES,
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildNodeInfo - build jarvismsg with NODE_INFO
func BuildNodeInfo(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string, destAddr string,
	ni *pb.NodeBaseInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_NODE_INFO,
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildFileData - build jarvismsg with TRANSFER_FILE
func BuildFileData(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string, destAddr string,
	fd *pb.FileData) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_TRANSFER_FILE,
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildRequestFile - build jarvismsg with REQUEST_FILE
func BuildRequestFile(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string, destAddr string,
	rf *pb.RequestFile) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REQUEST_FILE,
		Data: &pb.JarvisMsg_RequestFile{
			RequestFile: rf,
		},
	}

	err := SignJarvisMsg(privkey, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// BuildReplyRequestFile - build jarvismsg with REPLY_REQUEST_FILE
func BuildReplyRequestFile(privkey *jarviscrypto.PrivateKey, msgid int64, srcAddr string, destAddr string,
	fd *pb.FileData) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_REPLY_REQUEST_FILE,
		Data: &pb.JarvisMsg_File{
			File: fd,
		},
	}

	err := SignJarvisMsg(privkey, msg)
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
