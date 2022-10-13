package storage

import (
	"context"
	"fmt"

	"github.com/cestlascorpion/cuttlefish/utils"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	prefix string
	client *redis.Client
}

func NewRedis(ctx context.Context, conf *utils.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Network:  conf.Redis.Network,
		Addr:     conf.Redis.Addr,
		Username: conf.Redis.Username,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Database,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Errorf("redis ping err %+v", err)
		return nil, err
	}

	return &Redis{
		client: client,
		prefix: conf.Redis.KeyPrefix,
	}, nil
}

func (r *Redis) GetTentacle(ctx context.Context, id uint32) (map[string]string, error) {
	key := r.genKey(id)
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Errorf("redis hgetall %s err %+v", key, err)
		return nil, err
	}
	return result, nil
}

func (r *Redis) BatchGetTentacle(ctx context.Context, idList []uint32) (map[uint32]map[string]string, error) {
	resultList := make(map[uint32]map[string]string, len(idList))

	keyList := make([]string, 0, len(idList))
	for i := range idList {
		keyList = append(keyList, r.genKey(idList[i]))
	}

	pipeline := r.client.Pipeline()
	defer pipeline.Close()

	for i := range keyList {
		result, err := r.client.HGetAll(ctx, keyList[i]).Result()
		if err != nil {
			log.Errorf("redis hgetall %s err %+v", keyList[i], err)
			continue
		}
		resultList[idList[i]] = result
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Errorf("redis pipeliner exec err %+v", err)
		return nil, err
	}
	return resultList, nil
}

func (r *Redis) AddTentacle(ctx context.Context, id uint32, field string, value interface{}) error {
	key := r.genKey(id)
	_, err := r.client.HSet(ctx, key, field, value).Result()
	if err != nil {
		log.Errorf("redis hmset %s %s %+v err %+v", key, field, value, err)
		return err
	}
	return nil
}

func (r *Redis) AddMultiTentacle(ctx context.Context, id uint32, fields map[string]interface{}) error {
	key := r.genKey(id)
	_, err := r.client.HMSet(ctx, key, fields).Result()
	if err != nil {
		log.Errorf("redis hmset %s %+v err %+v", key, fields, err)
		return err
	}
	return nil
}

func (r *Redis) DelTentacle(ctx context.Context, id uint32, field string) (bool, error) {
	key := r.genKey(id)
	_, err := r.client.HDel(ctx, key, field).Result()
	if err != nil {
		log.Errorf("redis hdel %s %s err %+v", key, field, err)
		return false, err
	}

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("redis exists %s err %+v", key, err)
		return false, err
	}
	return exists > 0, nil
}

func (r *Redis) Close(ctx context.Context) error {
	return r.Close(ctx)
}

// ---------------------------------------------------------------------------------------------------------------------

func (r *Redis) genKey(id uint32) string {
	return fmt.Sprintf("%s_%d", r.prefix, id)
}
