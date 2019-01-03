package jarviscore

import (
	"context"
	"sync"

	pb "github.com/zhs007/jarviscore/proto"
)

// FuncReplyRequest - func reply request
type FuncReplyRequest func(ctx context.Context, jarvisnode JarvisNode, request *pb.JarvisMsg, reply *pb.JarvisMsg) (bool, error)

type requestData4Node struct {
	request   *pb.JarvisMsg
	funcReply FuncReplyRequest
}

type requestNodeData struct {
	mapData sync.Map
}

func (rnd *requestNodeData) addRequestData(request *pb.JarvisMsg, funcReply FuncReplyRequest) error {
	_, ok := rnd.mapData.Load(request.MsgID)
	if ok {
		return ErrDuplicateMsgID
	}

	rnd.mapData.Store(request.MsgID, &requestData4Node{
		request:   request,
		funcReply: funcReply,
	})

	return nil
}

func (rnd *requestNodeData) onReplyRequest(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error {
	request, ok := rnd.mapData.Load(msg.ReplyMsgID)
	if ok {
		rd4n := request.(*requestData4Node)
		if rd4n == nil {
			return ErrInvalidRequestData4Node
		}

		isfinished, err := rd4n.funcReply(ctx, jarvisnode, rd4n.request, msg)
		if err != nil {
			return err
		}

		if isfinished {
			rnd.mapData.Delete(msg.ReplyMsgID)
		}
	}

	return nil
}

type requestMgr struct {
	mapNode sync.Map
}

func (mgr *requestMgr) addRequestData(request *pb.JarvisMsg, funcReply FuncReplyRequest) error {
	dat, ok := mgr.mapNode.Load(request.DestAddr)
	if ok {
		rnd := dat.(*requestNodeData)
		if rnd == nil {
			return ErrInvalidRequestNodeData
		}

		return rnd.addRequestData(request, funcReply)
	}

	rnd := &requestNodeData{}
	mgr.mapNode.Store(request.DestAddr, rnd)

	return rnd.addRequestData(request, funcReply)
}

func (mgr *requestMgr) onReplyRequest(ctx context.Context, jarvisnode JarvisNode, msg *pb.JarvisMsg) error {
	dat, ok := mgr.mapNode.Load(msg.SrcAddr)
	if ok {
		rnd := dat.(*requestNodeData)
		if rnd == nil {
			return ErrInvalidRequestNodeData
		}

		return rnd.onReplyRequest(ctx, jarvisnode, msg)
	}

	return nil
}
