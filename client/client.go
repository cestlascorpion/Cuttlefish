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

func (c *Client) GetTentacle(ctx context.Context, uid uint32) ([]*pb.Tentacle, error) {
	resp, err := c.CuttlefishClient.GetTentacle(ctx, &pb.GetTentacleReq{
		Uid: uid,
	})
	if err != nil {
		return nil, err
	}
	return resp.Info, nil
}

func (c *Client) BatchGetTentacle(ctx context.Context, uidList []uint32) (map[uint32]*pb.TentacleList, error) {
	resp, err := c.CuttlefishClient.BatchGetTentacle(ctx, &pb.BatchGetTentacleReq{
		UidList: uidList,
	})
	if err != nil {
		return nil, err
	}
	return resp.InfoList, nil
}

func (c *Client) SetTentacle(ctx context.Context, uid uint32, key string, val map[string]string) (bool, error) {
	resp, err := c.CuttlefishClient.SetTentacle(ctx, &pb.SetTentacleReq{
		Info: &pb.Tentacle{
			Uid: uid,
			Key: key,
			Val: val,
		},
	})
	if err != nil {
		return false, err
	}
	return resp.Online, nil
}

func (c *Client) BatchSetTentacle(ctx context.Context) (map[uint32]bool, error) {
	// TODO:
	return nil, nil
}

func (c *Client) DelTentacle(ctx context.Context, uid uint32, key string) (bool, error) {
	resp, err := c.CuttlefishClient.DelTentacle(ctx, &pb.DelTentacleReq{
		Info: &pb.Tentacle{
			Uid: uid,
			Key: key,
		},
	})
	if err != nil {
		return false, err
	}
	return resp.Offline, nil
}

func (c *Client) BatchDelTentacle(ctx context.Context) (map[uint32]string, error) {
	// TODO:
	return nil, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.cc.Close()
}
