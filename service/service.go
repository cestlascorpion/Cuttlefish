package service

import (
	"context"

	pb "github.com/cestlascorpion/cuttlefish/proto"
	"github.com/cestlascorpion/cuttlefish/storage"
	"github.com/cestlascorpion/cuttlefish/utils"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	*pb.UnimplementedCuttlefishServer
	dao *storage.Redis
}

func NewServer(ctx context.Context, conf *utils.Config) (*Server, error) {
	if conf == nil {
		log.Errorf("conf is empty")
		return nil, utils.ErrInvalidParameter
	}

	d, err := storage.NewRedis(ctx, conf)
	if err != nil {
		log.Errorf("new redis err %+v", err)
		return nil, err
	}

	return &Server{
		dao: d,
	}, nil
}

func (s *Server) GetTentacle(ctx context.Context, in *pb.GetTentacleReq) (*pb.GetTentacleResp, error) {
	log.Debugf("GetTentacle in %+v", in)

	out := &pb.GetTentacleResp{
		// TODO:
	}
	// TODO:
	return out, nil
}

func (s *Server) BatchGetTentacle(ctx context.Context, in *pb.BatchGetTentacleReq) (*pb.BatchGetTentacleResp, error) {
	log.Debugf("BatchGetTentacle in %+v", in)

	out := &pb.BatchGetTentacleResp{
		// TODO:
	}
	// TODO:
	return out, nil
}

func (s *Server) SetTentacle(ctx context.Context, in *pb.SetTentacleReq) (*pb.SetTentacleResp, error) {
	log.Debugf("SetTentacle in %+v", in)

	out := &pb.SetTentacleResp{
		// TODO:
	}
	// TODO:
	return out, nil
}

func (s *Server) BatchSetTentacle(ctx context.Context, in *pb.BatchSetTentacleReq) (*pb.BatchSetTentacleResp, error) {
	log.Debugf("BatchSetTentacle in %+v", in)

	out := &pb.BatchSetTentacleResp{
		// TODO:
	}
	// TODO:
	return out, nil
}

func (s *Server) DelTentacle(ctx context.Context, in *pb.DelTentacleReq) (*pb.DelTentacleResp, error) {
	log.Debugf("DelTentacle in %+v", in)

	out := &pb.DelTentacleResp{
		// TODO:
	}
	// TODO:
	return out, nil
}

func (s *Server) BatchDelTentacle(ctx context.Context, in *pb.BatchDelTentacleReq) (*pb.BatchDelTentacleResp, error) {
	log.Debugf("BatchDelTentacle in %+v", in)

	out := &pb.BatchDelTentacleResp{
		// TODO:
	}
	// TODO:
	return out, nil
}

func (s *Server) Close(ctx context.Context) error {
	return s.dao.Close(ctx)
}
