package client

import (
	"context"

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

func (c *Client) GetTentacle(ctx context.Context, id uint32) (*pb.TentacleList, error) {
	resp, err := c.CuttlefishClient.GetTentacle(ctx, &pb.GetTentacleReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (c *Client) BatchGetTentacle(ctx context.Context, idList []uint32) (map[uint32]*pb.TentacleList, error) {
	resp, err := c.CuttlefishClient.BatchGetTentacle(ctx, &pb.BatchGetTentacleReq{
		IdList: idList,
	})
	if err != nil {
		return nil, err
	}
	if resp.Result == nil {
		return nil, nil
	}
	return resp.Result, nil
}

func (c *Client) SetTentacle(ctx context.Context, id uint32, list []*pb.Tentacle) error {
	_, err := c.CuttlefishClient.SetTentacle(ctx, &pb.SetTentacleReq{
		Id: id,
		List: &pb.TentacleList{
			TentacleList: list,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DelTentacle(ctx context.Context, id uint32, list []*pb.Tentacle) (bool, error) {
	resp, err := c.CuttlefishClient.DelTentacle(ctx, &pb.DelTentacleReq{
		Id: id,
		List: &pb.TentacleList{
			TentacleList: list,
		},
	})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.cc.Close()
}
