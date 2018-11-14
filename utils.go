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
//		sign(destAddr + curTime + srcAddr + data)
//		for mul-language, all become string merge data
func buildSignBuf(msg *pb.JarvisMsg) ([]byte, error) {
	if msg.MsgType == pb.MSGTYPE_LOCAL_CONNECT_OTHER ||
		msg.MsgType == pb.MSGTYPE_CONNECT_NODE {

		ci := msg.GetConnInfo()
		if ci != nil {
			str := []byte(fmt.Sprintf("%v%v%v", msg.DestAddr, msg.CurTime, msg.SrcAddr))
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
			str := []byte(fmt.Sprintf("%v%v%v", msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(ni)
			if err != nil {
				return nil, err
			}

			return append(str[:], buf[:]...), nil
		}
	} else if msg.MsgType == pb.MSGTYPE_REQUEST_CTRL {

		ci := msg.GetCtrlInfo()
		if ci != nil {
			str := []byte(fmt.Sprintf("%v%v%v", msg.DestAddr, msg.CurTime, msg.SrcAddr))
			buf, err := proto.Marshal(ci)
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
	destAddr string, servaddr string, ci *pb.CtrlInfo) (*pb.JarvisMsg, error) {

	msg := &pb.JarvisMsg{
		MsgID:    msgid,
		CurTime:  time.Now().Unix(),
		SrcAddr:  srcAddr,
		MyAddr:   srcAddr,
		DestAddr: destAddr,
		MsgType:  pb.MSGTYPE_LOCAL_CONNECT_OTHER,
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