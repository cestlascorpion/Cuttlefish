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

func TestRedis_SetTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer redis.Close(context.Background())

	err = redis.SetTentacle(context.Background(), 1234, "rock", "online")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestRedis_SetlMultiTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer redis.Close(context.Background())

	err = redis.SetMultiTentacle(context.Background(), 1234, map[string]interface{}{
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

func TestRedis_GetTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer redis.Close(context.Background())

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
	defer redis.Close(context.Background())

	resList, err := redis.BatchGetTentacle(context.Background(), []uint32{1234})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(resList)
}

func TestRedis_DelTentacle(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer redis.Close(context.Background())

	exist, err := redis.DelTentacle(context.Background(), 1234, "rock")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)
}

func TestRedis_DelTentacle2(t *testing.T) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer redis.Close(context.Background())

	exist, err := redis.DelTentacle(context.Background(), 1234, "r", "o", "c", "k")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(exist)
}

func BenchmarkRedis_SetTentacle(b *testing.B) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		b.FailNow()
	}
	defer redis.Close(context.Background())
	log.SetLevel(log.InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := redis.SetTentacle(context.Background(), 1234, "rock", "online")
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

}

func BenchmarkRedis_GetTentacle(b *testing.B) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		b.FailNow()
	}
	defer redis.Close(context.Background())
	log.SetLevel(log.InfoLevel)

	err = redis.SetTentacle(context.Background(), 1234, "rock", "online")
	if err != nil {
		fmt.Println(err)
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := redis.GetTentacle(context.Background(), 1234)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func BenchmarkRedis_BatchGetTentacle(b *testing.B) {
	conf := utils.NewTestConfig()
	redis, err := NewRedis(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		b.FailNow()
	}
	defer redis.Close(context.Background())
	log.SetLevel(log.InfoLevel)

	batchSize := 5120
	idList := make([]uint32, 0, batchSize)
	for i := 0; i < batchSize; i++ {
		idList = append(idList, uint32(i))
		err := redis.SetTentacle(context.Background(), uint32(i), "rock", "online")
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := redis.BatchGetTentacle(context.Background(), idList)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
