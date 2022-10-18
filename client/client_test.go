package client

import (
	"context"
	"fmt"
	pb "github.com/cestlascorpion/cuttlefish/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

const (
	target = "127.0.0.1:8080"
)

func TestClient_SetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	err = client.SetTentacle(context.Background(), 1234, []*pb.Tentacle{
		{
			Key: &pb.TentacleKey{
				Longitude: 1,
				Latitude:  1,
				Sequence:  1,
			},
			Val: &pb.TentacleVal{
				Connected: true,
				Timestamp: time.Now().UnixMilli(),
			},
		},
		{
			Key: &pb.TentacleKey{
				Longitude: 1,
				Latitude:  1,
				Sequence:  2,
			},
			Val: &pb.TentacleVal{
				Connected: true,
				Timestamp: time.Now().UnixMilli(),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestClient_GetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.GetTentacle(context.Background(), 1234)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(result)
}

func TestClient_DelTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	exists, err := client.DelTentacle(context.Background(), 1234, []*pb.Tentacle{
		{
			Key: &pb.TentacleKey{
				Longitude: 1,
				Latitude:  1,
				Sequence:  1,
			},
		},
		{
			Key: &pb.TentacleKey{
				Longitude: 1,
				Latitude:  1,
				Sequence:  2,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exists)
}

func TestClient_BatchGetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.BatchGetTentacle(context.Background(), []uint32{1234})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(result)
}
