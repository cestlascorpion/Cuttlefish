package service

import (
	"context"
	"encoding/json"
	"time"

	pb "github.com/cestlascorpion/cuttlefish/proto"
	"github.com/cestlascorpion/cuttlefish/storage"
	"github.com/cestlascorpion/cuttlefish/utils"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	*pb.UnimplementedCuttlefishServer
	dao        *storage.Redis
	histChan   chan *histMeta
	histCancel []context.CancelFunc
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

	svr := &Server{
		dao:        d,
		histChan:   make(chan *histMeta, defaultChanBuffer),
		histCancel: make([]context.CancelFunc, 0),
	}

	for i := 0; i < defaultParallel; i++ {
		x, cancel := context.WithCancel(ctx)
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case m, ok := <-svr.histChan:
					if !ok {
						return
					}
					if m != nil {
						_ = svr.handleHistMeta(ctx, m)
					}
				}
			}
		}(x)
		svr.histCancel = append(svr.histCancel, cancel)
	}

	return svr, nil
}

func (s *Server) GetTentacle(ctx context.Context, in *pb.GetTentacleReq) (*pb.GetTentacleResp, error) {
	log.Debugf("GetTentacle in %+v", in)

	out := &pb.GetTentacleResp{
		TentacleList: make([]*pb.Tentacle, 0),
	}

	result, err := s.dao.GetTentacle(ctx, in.Uid)
	if err != nil {
		log.Errorf("GetTentacle %d err %+v", in.Uid, err)
		return out, err
	}

	for k, v := range result {
		out.TentacleList = append(out.TentacleList, &pb.Tentacle{
			Key: k,
			Val: v,
		})
	}
	return out, nil
}

func (s *Server) BatchGetTentacle(ctx context.Context, in *pb.BatchGetTentacleReq) (*pb.BatchGetTentacleResp, error) {
	log.Debugf("BatchGetTentacle in %+v", in)

	out := &pb.BatchGetTentacleResp{
		InfoList: make(map[uint32]*pb.TentacleInfo, len(in.UidList)),
	}
	if len(in.UidList) == 0 {
		log.Errorf("BatchGetTentacle empty uidList")
		return out, utils.ErrInvalidParameter
	}

	resultMap, err := s.dao.BatchGetTentacle(ctx, in.UidList)
	if err != nil {
		log.Errorf("BatchGetTentacle %v err %+v", in.UidList, err)
		return out, err
	}

	for id, result := range resultMap {
		infoList := make([]*pb.Tentacle, 0, len(result))
		for k, v := range result {
			infoList = append(infoList, &pb.Tentacle{
				Key: k,
				Val: v,
			})
		}
		out.InfoList[id] = &pb.TentacleInfo{
			TentacleList: infoList,
		}
		log.Debugf("BatchGetTentacle %d %+v", id, infoList)
	}
	return out, nil
}

func (s *Server) SetTentacle(ctx context.Context, in *pb.SetTentacleReq) (*pb.SetTentacleResp, error) {
	log.Debugf("SetTentacle in %+v", in)

	out := &pb.SetTentacleResp{
		Online: false,
	}
	if len(in.TentacleList) == 0 {
		log.Errorf("SetTentacle empty tentacleList")
		return out, utils.ErrInvalidParameter
	}

	exists, err := s.dao.SetTentacle(ctx, in.Id, in.TentacleList)
	if err != nil {
		log.Errorf("SetTentacle %d %+v err %+v", in.Id, in.TentacleList, err)
		return out, err
	}
	out.Online = !exists

	ts := time.Now()
	if !exists {
		s.histChan <- &histMeta{
			id:     in.Id,
			ts:     ts,
			status: true,
		}
	}

	return out, nil
}

func (s *Server) BatchSetTentacle(ctx context.Context, in *pb.BatchSetTentacleReq) (*pb.BatchSetTentacleResp, error) {
	log.Debugf("BatchSetTentacle in %+v", in)

	out := &pb.BatchSetTentacleResp{
		Result: make(map[uint32]bool, len(in.InfoList)),
	}
	if len(in.InfoList) == 0 {
		log.Errorf("BatchSetTentacle empty infoList")
		return out, utils.ErrInvalidParameter
	}

	existsMap, err := s.dao.BatchSetTentacle(ctx, in.InfoList)
	if err != nil {
		log.Errorf("BatchSetTentacle %+v err %+v", in.InfoList, err)
		return out, err
	}

	ts := time.Now()
	for id, exists := range existsMap {
		out.Result[id] = !exists
		if !exists {
			s.histChan <- &histMeta{
				id:     id,
				ts:     ts,
				status: true,
			}
		}
	}
	return out, nil
}

func (s *Server) DelTentacle(ctx context.Context, in *pb.DelTentacleReq) (*pb.DelTentacleResp, error) {
	log.Debugf("DelTentacle in %+v", in)

	out := &pb.DelTentacleResp{
		Offline: false,
	}
	if len(in.TentacleList) == 0 {
		log.Errorf("SetTentacle empty tentacleList")
		return out, utils.ErrInvalidParameter
	}

	exists, err := s.dao.DelTentacle(ctx, in.Id, in.TentacleList)
	if err != nil {
		log.Errorf("DelTentacle %d %+v err %+v", in.Id, in.TentacleList, err)
		return out, nil
	}
	out.Offline = !exists
	if !exists {
		s.histChan <- &histMeta{
			id:     in.Id,
			ts:     time.Now(),
			status: false,
		}
	}
	return out, nil
}

func (s *Server) BatchDelTentacle(ctx context.Context, in *pb.BatchDelTentacleReq) (*pb.BatchDelTentacleResp, error) {
	log.Debugf("BatchDelTentacle in %+v", in)

	out := &pb.BatchDelTentacleResp{
		Result: make(map[uint32]bool, len(in.InfoList)),
	}
	if len(in.InfoList) == 0 {
		log.Errorf("BatchDelTentacle empty infoList")
		return out, utils.ErrInvalidParameter
	}

	resultMap, err := s.dao.BatchDelTentacle(ctx, in.InfoList)
	if err != nil {
		log.Errorf("BatchDelTentacle %+v err %+v", in.InfoList, err)
		return out, err
	}

	for id, exists := range resultMap {
		out.Result[id] = !exists
		if !exists {
			s.histChan <- &histMeta{
				id:     id,
				ts:     time.Now(),
				status: false,
			}
		}
	}
	return out, nil
}

func (s *Server) GetTentacleHistory(ctx context.Context, in *pb.GetTentacleHistoryReq) (*pb.GetTentacleHistoryResp, error) {
	log.Debugf("GetTentacleHistory in %+v", in)

	out := &pb.GetTentacleHistoryResp{
		InfoList: make([]*pb.HistoryInfo, 0),
	}
	result, err := s.dao.GetTentacleHistory(ctx, in.Id, time.UnixMilli(in.From), time.UnixMilli(in.To))
	if err != nil {
		log.Errorf("GetTentacleHistory %d err %+v", in.Id, err)
		return out, err
	}

	for i := range result {
		info := &pb.HistoryInfo{}
		err := json.Unmarshal([]byte(result[i]), info)
		if err != nil {
			log.Errorf("proto unmarshal err %+v", err)
			return out, err
		}
		out.InfoList = append(out.InfoList, info)
	}
	return out, nil
}

func (s *Server) Close(ctx context.Context) error {
	close(s.histChan)
	log.Info("histMeta chan closed")
	for i := range s.histCancel {
		s.histCancel[i]()
	}
	log.Info("histMeta routine canceled")
	for m := range s.histChan {
		if m != nil {
			_ = s.handleHistMeta(ctx, m)
		}
	}
	log.Info("histMeta chan clear")
	return s.dao.Close(ctx)
}

// ---------------------------------------------------------------------------------------------------------------------

const (
	defaultChanBuffer = 4096
	defaultParallel   = 10
)

type histMeta struct {
	id     uint32
	ts     time.Time
	status bool
}

func (s *Server) handleHistMeta(ctx context.Context, data *histMeta) error {
	info := &pb.HistoryInfo{
		St: data.status,
		Ts: data.ts.UnixMilli(),
	}

	bs, err := json.Marshal(info)
	if err != nil {
		log.Errorf("json marshal err %+v", err)
		return err
	}
	record := string(bs)
	err = s.dao.SetTentacleHistory(ctx, data.id, record, data.ts)
	if err != nil {
		log.Errorf("SetTentacleHistory %d %s %+v err %+v", data.id, record, data.ts, err)
		return err
	}
	return nil
}
