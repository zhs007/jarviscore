package jarviscore

import (
	"context"

	pb "github.com/zhs007/jarviscore/proto"
)

// server is used to implement jarviscorepb.JarvisCoreServ.
type server struct{}

// Join implements jarviscorepb.JarvisCoreServ
func (s *server) Join(ctx context.Context, in *pb.Join) (*pb.ReplyJoin, error) {
	return &pb.ReplyJoin{Code: pb.CODE_OK}, nil
}
