package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/cestlascorpion/cuttlefish/utils"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	prefix string
	expire time.Duration
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
		expire: time.Second * time.Duration(conf.Redis.KeyExpire),
	}, nil
}

func (r *Redis) SetTentacle(ctx context.Context, id uint32, field string, value interface{}) error {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	key := r.genKey(id)
	hSetCmd := pipe.HSet(ctx, key, field, value)
	expireCmd := pipe.Expire(ctx, key, r.expire)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return err
	}

	err = hSetCmd.Err()
	if err != nil {
		log.Errorf("hset %s %s %v err %+v", key, field, value, err)
		return err
	}

	err = expireCmd.Err()
	if err != nil {
		log.Errorf("expire %s err %+v", key, err)
		return err
	}

	return nil
}

func (r *Redis) SetMultiTentacle(ctx context.Context, id uint32, fields map[string]interface{}) error {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	key := r.genKey(id)
	hmSetCmd := pipe.HMSet(ctx, key, fields)
	expireCmd := pipe.Expire(ctx, key, r.expire)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return err
	}

	err = hmSetCmd.Err()
	if err != nil {
		log.Errorf("hmset %s %v err %+v", key, fields, err)
		return err
	}

	err = expireCmd.Err()
	if err != nil {
		log.Errorf("expire %s err %+v", key, err)
		return err
	}

	return nil
}

func (r *Redis) GetTentacle(ctx context.Context, id uint32) (map[string]string, error) {
	key := r.genKey(id)
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Errorf("hgetall %s err %+v", key, err)
		return nil, err
	}
	return result, nil
}

func (r *Redis) BatchGetTentacle(ctx context.Context, idList []uint32) (map[uint32]map[string]string, error) {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	result := make([]*redis.StringStringMapCmd, 0)
	for i := range idList {
		result = append(result, pipe.HGetAll(ctx, r.genKey(idList[i])))
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return nil, err
	}

	resultList := make(map[uint32]map[string]string, len(idList))
	for i := range result {
		result, err := result[i].Result()
		if err != nil {
			log.Warnf("hgetall %s err %+v", r.genKey(idList[i]), err)
			continue
		}
		resultList[idList[i]] = result
	}
	return resultList, nil
}

func (r *Redis) DelTentacle(ctx context.Context, id uint32, fields ...string) (bool, error) {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	key := r.genKey(id)
	hDelCmd := r.client.HDel(ctx, key, fields...)
	existsCmd := r.client.Exists(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return false, err
	}

	err = hDelCmd.Err()
	if err != nil {
		log.Errorf("hdel %s %v err %+v", key, fields, err)
		return false, err
	}

	err = existsCmd.Err()
	if err != nil {
		log.Errorf("exists %s err %+v", key, err)
		return false, err
	}

	return existsCmd.Val() > 0, nil
}

func (r *Redis) Close(ctx context.Context) error {
	return r.client.Close()
}

// ---------------------------------------------------------------------------------------------------------------------

func (r *Redis) genKey(id uint32) string {
	return fmt.Sprintf("%s_%d", r.prefix, id)
}
