package server

import (
	"context"
	"github.com/Rellum/inventive_weave/svc/creators/activity"
	pb "github.com/Rellum/inventive_weave/svc/creators/creatorspb"
	"google.golang.org/grpc"
)

func RegisterServer(s *grpc.Server) {
	pb.RegisterCreatorsServer(s, new(server))
}

type server struct {
	pb.UnimplementedCreatorsServer
}

func (s server) MostActiveCreators(ctx context.Context, req *pb.MostActiveCreatorsReq) (*pb.MostActiveCreatorsRes, error) {
	return &pb.MostActiveCreatorsRes{
		CreatorStats: activity.MostActive(req.Creators, req.Products),
	}, nil
}
