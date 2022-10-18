package service

import (
	"context"

	pb "github.com/cestlascorpion/cuttlefish/proto"
	"github.com/cestlascorpion/cuttlefish/storage"
	"github.com/cestlascorpion/cuttlefish/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
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
		Result: &pb.TentacleList{
			TentacleList: make([]*pb.Tentacle, 0),
		},
	}
	result, err := s.dao.GetTentacle(ctx, in.Id)
	if err != nil {
		log.Errorf("GetTentacle err %+v", err)
		return out, err
	}
	for k, v := range result {
		key, val, err := unmarshal(k, v)
		if err != nil {
			log.Warnf("unmarshal %s %s err %+v", k, v, err)
			continue
		}
		out.Result.TentacleList = append(out.Result.TentacleList, &pb.Tentacle{
			Key: key,
			Val: val,
		})
	}
	return out, nil
}

func (s *Server) BatchGetTentacle(ctx context.Context, in *pb.BatchGetTentacleReq) (*pb.BatchGetTentacleResp, error) {
	log.Debugf("BatchGetTentacle in %+v", in)

	out := &pb.BatchGetTentacleResp{
		Result: make(map[uint32]*pb.TentacleList, len(in.IdList)),
	}
	if len(in.IdList) == 0 {
		log.Errorf("empty id list")
		return out, utils.ErrInvalidParameter
	}

	result, err := s.dao.BatchGetTentacle(ctx, in.IdList)
	if err != nil {
		log.Errorf("BatchGetTentacle err %+v", err)
		return out, err
	}
	for id, res := range result {
		list := &pb.TentacleList{
			TentacleList: make([]*pb.Tentacle, 0),
		}
		for k, v := range res {
			key, val, err := unmarshal(k, v)
			if err != nil {
				log.Warnf("unmarshal %s %s err %+v", k, v, err)
				continue
			}
			list.TentacleList = append(list.TentacleList, &pb.Tentacle{
				Key: key,
				Val: val,
			})
		}
		if len(list.TentacleList) > 0 {
			out.Result[id] = list
		}
	}
	return out, nil
}

func (s *Server) SetTentacle(ctx context.Context, in *pb.SetTentacleReq) (*pb.SetTentacleResp, error) {
	log.Debugf("SetTentacle in %+v", in)

	out := &pb.SetTentacleResp{}
	if in.List == nil || len(in.List.TentacleList) == 0 {
		log.Errorf("empty tentacle list")
		return out, utils.ErrInvalidParameter
	}

	if len(in.List.TentacleList) == 1 {
		k, v, err := marshal(in.List.TentacleList[0].Key, in.List.TentacleList[0].Val)
		if err != nil {
			log.Errorf("marshal %v %v err %+v", in.List.TentacleList[0].Key, in.List.TentacleList[0].Val, err)
			return out, err
		}
		err = s.dao.SetTentacle(ctx, in.Id, k, v)
		if err != nil {
			log.Errorf("SetTentacle err %+v", err)
			return out, err
		}
		return out, nil
	}

	fields := make(map[string]interface{}, len(in.List.TentacleList))
	for i := range in.List.TentacleList {
		k, v, err := marshal(in.List.TentacleList[i].Key, in.List.TentacleList[i].Val)
		if err != nil {
			log.Errorf("marshal %v %v err %+v", in.List.TentacleList[i].Key, in.List.TentacleList[i].Val, err)
			return out, err
		}
		fields[k] = v
	}
	err := s.dao.SetMultiTentacle(ctx, in.Id, fields)
	if err != nil {
		log.Errorf("SetMultiTentacle err %+v", err)
		return out, err
	}
	return out, nil
}

func (s *Server) DelTentacle(ctx context.Context, in *pb.DelTentacleReq) (*pb.DelTentacleResp, error) {
	log.Debugf("DelTentacle in %+v", in)

	out := &pb.DelTentacleResp{}
	if in.List == nil || len(in.List.TentacleList) == 0 {
		log.Errorf("empty tentacle list")
		return out, utils.ErrInvalidParameter
	}

	fields := make([]string, 0, len(in.List.TentacleList))
	for i := range in.List.TentacleList {
		k, err := marshalKey(in.List.TentacleList[i].Key)
		if err != nil {
			log.Errorf("marshalKey %v err %+v", in.List.TentacleList[i].Key, err)
			return out, err
		}
		fields = append(fields, k)
	}
	exists, err := s.dao.DelTentacle(ctx, in.Id, fields...)
	if err != nil {
		log.Errorf("DelTentacle err %+v", err)
		return out, err
	}
	out.Exists = exists
	return out, nil
}

func (s *Server) Close(ctx context.Context) error {
	return s.dao.Close(ctx)
}

// ---------------------------------------------------------------------------------------------------------------------

func marshal(key *pb.TentacleKey, val *pb.TentacleVal) (string, string, error) {
	if key == nil || val == nil {
		return "", "", utils.ErrInvalidParameter
	}
	k, err := proto.Marshal(key)
	if err != nil {
		return "", "", err
	}
	v, err := proto.Marshal(val)
	if err != nil {
		return "", "", err
	}
	return string(k), string(v), nil
}

func marshalKey(key *pb.TentacleKey) (string, error) {
	if key == nil {
		return "", utils.ErrInvalidParameter
	}
	k, err := proto.Marshal(key)
	if err != nil {
		return "", err
	}
	return string(k), nil
}

func unmarshal(k, v string) (*pb.TentacleKey, *pb.TentacleVal, error) {
	key, val := &pb.TentacleKey{}, &pb.TentacleVal{}
	err := proto.Unmarshal([]byte(k), key)
	if err != nil {
		return nil, nil, err
	}
	err = proto.Unmarshal([]byte(v), val)
	if err != nil {
		return nil, nil, err
	}
	return key, val, nil
}
