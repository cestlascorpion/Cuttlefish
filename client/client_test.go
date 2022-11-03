package client

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	pb "github.com/cestlascorpion/cuttlefish/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	target = "127.0.0.1:8080"
)

var (
	egK = &pb.EgK{
		Proxy: "127.0.0.1:80",
		SeqId: "0XFF",
	}
	egV = &pb.EgV{
		Ts:   time.Now().UnixMilli(),
		Desc: "pc",
	}

	egKey string
	egVal string
)

func init() {
	log.SetLevel(log.DebugLevel)
	k, err := json.Marshal(egK)
	if err != nil {
		fmt.Println("marshal err", err)
	} else {
		egKey = string(k)
	}
	v, err := json.Marshal(egV)
	if err != nil {
		fmt.Println("marshal err", err)
	} else {
		egVal = string(v)
	}
}

func TestClient_SetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	online, err := client.SetTentacle(context.Background(), 1234, "key", "val")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(online)
}

func TestClient_PeekTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	exists, err := client.PeekTentacle(context.Background(), 1234)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exists)
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

	offline, err := client.DelTentacle(context.Background(), 1234, "key")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(offline)
}

func TestClient_SetMultiTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	online, err := client.SetMultiTentacle(context.Background(), 1234, []*pb.Tentacle{
		{Key: egKey, Val: egVal},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(online)
}

func TestClient_GetTentacle2(t *testing.T) {
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
	for i := range result {
		key := &pb.EgK{}
		val := &pb.EgV{}
		err = json.Unmarshal([]byte(result[i].Key), key)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		err = json.Unmarshal([]byte(result[i].Val), val)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		fmt.Println(key, val)
	}
}

func TestClient_DelMultiTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	offline, err := client.DelMultiTentacle(context.Background(), 1234, []*pb.Tentacle{
		{Key: egKey},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(offline)
}

func TestClient_BatchSetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.BatchSetTentacle(context.Background(), map[uint32]*pb.TentacleInfo{
		0: {
			TentacleList: []*pb.Tentacle{
				{Key: "k", Val: "k"},
				{Key: "e", Val: "e"},
				{Key: "y", Val: "y"},
			},
		},
		1: {
			TentacleList: []*pb.Tentacle{
				{Key: "k", Val: "k"},
				{Key: "e", Val: "e"},
				{Key: "y", Val: "y"},
			},
		},
		2: {
			TentacleList: []*pb.Tentacle{
				{Key: "k", Val: "k"},
				{Key: "e", Val: "e"},
				{Key: "y", Val: "y"},
			},
		},
		3: {
			TentacleList: []*pb.Tentacle{
				{Key: "k", Val: "k"},
				{Key: "e", Val: "e"},
				{Key: "y", Val: "y"},
			},
		},
		4: {
			TentacleList: []*pb.Tentacle{
				{Key: "k", Val: "k"},
				{Key: "e", Val: "e"},
				{Key: "y", Val: "y"},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for k, v := range result {
		fmt.Println(k, v)
	}
}

func TestClient_BatchPeekTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.BatchPeekTentacle(context.Background(), []uint32{0, 1, 2, 3, 4})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(result)
}

func TestClient_BatchGetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.BatchGetTentacle(context.Background(), []uint32{0, 1, 2, 3, 4})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for k, v := range result {
		fmt.Println(k, v)
	}
}

func TestClient_BatchDelTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.BatchDelTentacle(context.Background(), map[uint32]*pb.TentacleInfo{
		0: {
			TentacleList: []*pb.Tentacle{
				{Key: "k"},
				{Key: "e"},
				{Key: "y"},
			},
		},
		1: {
			TentacleList: []*pb.Tentacle{
				{Key: "k"},
				{Key: "e"},
				{Key: "y"},
			},
		},
		2: {
			TentacleList: []*pb.Tentacle{
				{Key: "k"},
				{Key: "e"},
				{Key: "y"},
			},
		},
		3: {
			TentacleList: []*pb.Tentacle{
				{Key: "k"},
				{Key: "e"},
				{Key: "y"},
			},
		},
		4: {
			TentacleList: []*pb.Tentacle{
				{Key: "k"},
				{Key: "e"},
				{Key: "y"},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for k, v := range result {
		fmt.Println(k, v)
	}
}

func TestClient_GetTentacleHistory(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	result, err := client.GetTentacleHistory(context.Background(), 1234, time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for i := range result {
		fmt.Println(result[i])
	}
}
