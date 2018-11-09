package jarviscore

import (
	"fmt"
	"math/big"

	"github.com/zhs007/jarviscore/crypto"
	pb "github.com/zhs007/jarviscore/proto"
)

// buildSignBuf - build sign buf
//		sign(destAddr + curTime + data + srcAddr)
func buildSignBuf(msg *pb.JarvisMsg) []byte {
	if msg.MsgType == pb.MSGTYPE_JOIN {
		return []byte(fmt.Sprintf("%v%v%+v%v", msg.DestAddr, msg.CurTime, msg.GetJoinInfo(), msg.SrcAddr))
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
