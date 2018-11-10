package jarviscore

import (
	"fmt"
	"math/big"
	"time"

	"github.com/zhs007/jarviscore/crypto"
	pb "github.com/zhs007/jarviscore/proto"
)

// buildSignBuf - build sign buf
//		sign(destAddr + curTime + data + srcAddr)
//		for mul-language, all become string merge data
func buildSignBuf(msg *pb.JarvisMsg) []byte {
	if msg.MsgType == pb.MSGTYPE_CONNECT_NODE || msg.MsgType == pb.MSGTYPE_NODE_INFO || msg.MsgType == pb.MSGTYPE_REPLY_CONNECT {
		ni := msg.GetNodeInfo()
		if ni != nil {
			return []byte(fmt.Sprintf("%v%v%v%v%v%v", msg.DestAddr, msg.CurTime, ni.ServAddr, ni.Addr, ni.Name, msg.SrcAddr))
		}
	}

	return nil
}

// SignJarvisMsg - sign JarvisMsg
func SignJarvisMsg(privkey *jarviscrypto.PrivateKey, msg *pb.JarvisMsg) error {
	buf := buildSignBuf(msg)

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

	buf := buildSignBuf(msg)

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
func BuildConnNode(msgid int64, srcAddr string, destAddr string, ni *pb.NodeBaseInfo) *pb.JarvisMsg {
	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_CONNECT_NODE,
		Data: &pb.JarvisMsg_NodeInfo{
			NodeInfo: ni,
		},
	}

	return msg
}

// BuildReplyConn - build jarvismsg with REPLY_CONNECT
func BuildReplyConn(msgid int64, srcAddr string, destAddr string, ni *pb.NodeBaseInfo) *pb.JarvisMsg {
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

	return msg
}
