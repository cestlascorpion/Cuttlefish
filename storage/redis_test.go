package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/cestlascorpion/cuttlefish/proto"
	"github.com/cestlascorpion/cuttlefish/utils"
	log "github.com/sirupsen/logrus"
)

var dao *Redis

func init() {
	log.SetLevel(log.DebugLevel)

	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	dao = redis
}

func TestRedis_SetTentacle(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	exists, err := dao.SetTentacle(context.Background(), 1234, []*pb.Tentacle{
		{Key: "key", Val: "value"},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exists)
}

func TestRedis_GetTentacle(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	result, err := dao.GetTentacle(context.Background(), 1234)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(result)
}

func TestRedis_DelTentacle(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	exists, err := dao.DelTentacle(context.Background(), 1234, []*pb.Tentacle{
		{Key: "key"},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exists)
}

func TestRedis_SetTentacle2(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	exists, err := dao.SetTentacle(context.Background(), 1234, []*pb.Tentacle{
		{Key: "k1", Val: "v1"},
		{Key: "k2", Val: "v2"},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exists)
}

func TestRedis_GetTentacle2(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	result, err := dao.GetTentacle(context.Background(), 1234)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(result)
}

func TestRedis_DelTentacle2(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	exists, err := dao.DelTentacle(context.Background(), 1234, []*pb.Tentacle{
		{Key: "k1"},
		{Key: "k2"},
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exists)
}

func TestRedis_BatchSetTentacle(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	parameter := make(map[uint32]*pb.TentacleInfo)
	for i := uint32(0); i < 5; i++ {
		parameter[i] = &pb.TentacleInfo{
			TentacleList: []*pb.Tentacle{
				{Key: "k", Val: "k"},
				{Key: "e", Val: "e"},
				{Key: "y", Val: "y"},
			},
		}
	}
	existsMap, err := dao.BatchSetTentacle(context.Background(), parameter)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(existsMap)
}

func TestRedis_BatchGetTentacle(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	idList := make([]uint32, 0)
	for i := uint32(0); i < 5; i++ {
		idList = append(idList, i)
	}
	resultMap, err := dao.BatchGetTentacle(context.Background(), idList)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(resultMap)
}

func TestRedis_BatchDelTentacle(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	parameter := make(map[uint32]*pb.TentacleInfo)
	for i := uint32(0); i < 5; i++ {
		parameter[i] = &pb.TentacleInfo{
			TentacleList: []*pb.Tentacle{
				{Key: "k"},
				{Key: "e"},
				{Key: "y"},
			},
		}
	}
	existsMap, err := dao.BatchDelTentacle(context.Background(), parameter)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(existsMap)
}

func TestRedis_SetTentacleHistory(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	ts := time.Now().AddDate(0, 0, -1)
	err := dao.SetTentacleHistory(context.Background(), 1234, "online @"+ts.String(), ts)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	time.Sleep(time.Second)

	ts = time.Now().AddDate(0, 0, -1)
	err = dao.SetTentacleHistory(context.Background(), 1234, "offline @"+ts.String(), ts)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestRedis_GetTentacleHistory(t *testing.T) {
	if dao == nil {
		fmt.Println("init dao failed")
		return
	}

	result, err := dao.GetTentacleHistory(context.Background(), 1234, time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(result)
}
