package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/cestlascorpion/cuttlefish/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRedis_GetTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	res, err := redis.GetTentacle(context.Background(), 1234)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(res)
}

func TestRedis_BatchGetTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	resList, err := redis.BatchGetTentacle(context.Background(), []uint32{1, 2, 3, 4, 1234})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(resList)
}

func TestRedis_AddTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = redis.AddTentacle(context.Background(), 1234, "rock", "online")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = redis.AddTentacle(context.Background(), 1, "rock", "online")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = redis.AddTentacle(context.Background(), 2, "rock", "online")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = redis.AddTentacle(context.Background(), 3, "rock", "online")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = redis.AddTentacle(context.Background(), 4, "rock", "online")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestRedis_AddMultiTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = redis.AddMultiTentacle(context.Background(), 1234, map[string]interface{}{
		"r": "online",
		"o": "online",
		"c": "online",
		"k": "online",
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestRedis_DelTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	exist, err := redis.DelTentacle(context.Background(), 1234, "rock")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 1234, "r")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 1234, "o")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 1234, "c")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 1234, "k")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 1, "rock")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 2, "rock")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 3, "rock")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)

	exist, err = redis.DelTentacle(context.Background(), 4, "rock")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)
}
