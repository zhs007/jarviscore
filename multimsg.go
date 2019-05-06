package jarviscore

import (
	"github.com/golang/protobuf/proto"
	"github.com/zhs007/jarviscore/basedef"
	pb "github.com/zhs007/jarviscore/proto"
)

func makeMultiMsg(jarvisnode JarvisNode, msg *pb.JarvisMsg) ([]*pb.JarvisMsg, error) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	bl := len(buf)
	if bl > basedef.BigMsgLength {
		var msgs []*pb.JarvisMsg

		curstart := int64(0)
		curlength := int64(basedef.BigMsgLength)

		for curstart < int64(bl) {
			if curstart+curlength >= int64(bl) {

				cm := &pb.MultiMsgData{
					TotalMsgLength: int64(bl),
					CurMsgIndex:    int32(curstart),
					Buf:            buf[curstart:bl],
					CurMD5:         GetMD5String(buf[curstart:bl]),
					TotalMD5:       GetMD5String(buf),
				}

				curjm, err := BuildMultiMsgData(jarvisnode, msg.DestAddr, msg.ReplyMsgID, cm)
				if err != nil {
					return nil, err
				}

				msgs = append(msgs, curjm)

				break
			}

			cm := &pb.MultiMsgData{
				TotalMsgLength: int64(bl),
				CurMsgIndex:    int32(curstart),
				Buf:            buf[curstart:(curstart + curlength)],
				CurMD5:         GetMD5String(buf[curstart:(curstart + curlength)]),
			}

			curjm, err := BuildMultiMsgData(jarvisnode, msg.DestAddr, msg.ReplyMsgID, cm)
			if err != nil {
				return nil, err
			}

			msgs = append(msgs, curjm)

			curstart = curstart + curlength
		}

		return msgs, nil
	}

	return nil, nil
}
