package client

import (
	"context"
	"time"

	pb "github.com/cestlascorpion/cuttlefish/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Client struct {
	pb.CuttlefishClient
	cc *grpc.ClientConn
}

func NewClient(target string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Errorf("grpc dial %s err %+v", target, err)
		return nil, err
	}
	return &Client{
		CuttlefishClient: pb.NewCuttlefishClient(conn),
		cc:               conn,
	}, nil
}

func (c *Client) GetTentacle(ctx context.Context, uid uint32) ([]*pb.Tentacle, error) {
	resp, err := c.CuttlefishClient.GetTentacle(ctx, &pb.GetTentacleReq{
		Uid: uid,
	})
	if err != nil {
		return nil, err
	}
	return resp.TentacleList, nil
}

func (c *Client) BatchGetTentacle(ctx context.Context, uidList []uint32) (map[uint32]*pb.TentacleInfo, error) {
	resp, err := c.CuttlefishClient.BatchGetTentacle(ctx, &pb.BatchGetTentacleReq{
		UidList: uidList,
	})
	if err != nil {
		return nil, err
	}
	return resp.InfoList, nil
}

func (c *Client) PeekTentacle(ctx context.Context, uid uint32) (bool, error) {
	resp, err := c.CuttlefishClient.PeekTentacle(ctx, &pb.PeekTentacleReq{
		Uid: uid,
	})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

func (c *Client) BatchPeekTentacle(ctx context.Context, uidList []uint32) (map[uint32]bool, error) {
	resp, err := c.CuttlefishClient.BatchPeekTentacle(ctx, &pb.BatchPeekTentacleReq{
		UidList: uidList,
	})
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (c *Client) SetTentacle(ctx context.Context, uid uint32, key string, val string) (bool, error) {
	resp, err := c.CuttlefishClient.SetTentacle(ctx, &pb.SetTentacleReq{
		Id: uid,
		TentacleList: []*pb.Tentacle{
			{
				Key: key,
				Val: val,
			},
		},
	})
	if err != nil {
		return false, err
	}
	return resp.Online, nil
}

func (c *Client) SetMultiTentacle(ctx context.Context, uid uint32, infoList []*pb.Tentacle) (bool, error) {
	resp, err := c.CuttlefishClient.SetTentacle(ctx, &pb.SetTentacleReq{
		Id:           uid,
		TentacleList: infoList,
	})
	if err != nil {
		return false, err
	}
	return resp.Online, nil
}

func (c *Client) BatchSetTentacle(ctx context.Context, infoList map[uint32]*pb.TentacleInfo) (map[uint32]bool, error) {
	resp, err := c.CuttlefishClient.BatchSetTentacle(ctx, &pb.BatchSetTentacleReq{
		InfoList: infoList,
	})
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (c *Client) DelTentacle(ctx context.Context, uid uint32, key string) (bool, error) {
	resp, err := c.CuttlefishClient.DelTentacle(ctx, &pb.DelTentacleReq{
		Id: uid,
		TentacleList: []*pb.Tentacle{
			{
				Key: key,
			},
		},
	})
	if err != nil {
		return false, err
	}
	return resp.Offline, nil
}

func (c *Client) DelMultiTentacle(ctx context.Context, uid uint32, infoList []*pb.Tentacle) (bool, error) {
	resp, err := c.CuttlefishClient.DelTentacle(ctx, &pb.DelTentacleReq{
		Id:           uid,
		TentacleList: infoList,
	})
	if err != nil {
		return false, err
	}
	return resp.Offline, nil
}

func (c *Client) BatchDelTentacle(ctx context.Context, infoList map[uint32]*pb.TentacleInfo) (map[uint32]bool, error) {
	resp, err := c.CuttlefishClient.BatchDelTentacle(ctx, &pb.BatchDelTentacleReq{
		InfoList: infoList,
	})
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (c *Client) GetTentacleHistory(ctx context.Context, id uint32, from, to time.Time) ([]*pb.HistoryInfo, error) {
	resp, err := c.CuttlefishClient.GetTentacleHistory(ctx, &pb.GetTentacleHistoryReq{
		Id:   id,
		From: from.UnixMilli(),
		To:   to.UnixMilli(),
	})
	if err != nil {
		return nil, err
	}
	return resp.InfoList, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.cc.Close()
}
