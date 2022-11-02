package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/cestlascorpion/cuttlefish/proto"
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

func (r *Redis) SetTentacle(ctx context.Context, id uint32, infoList []*pb.Tentacle) (bool, error) {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	key := r.genUserKey(id)
	values := make([]interface{}, 0, len(infoList)*2)
	for i := range infoList {
		values = append(values, infoList[i].Key, infoList[i].Val)
	}

	existsCmd := pipe.Exists(ctx, key)
	hSetCmd := pipe.HSet(ctx, key, values...)
	expireCmd := pipe.Expire(ctx, key, r.expire)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return false, err
	}

	exists, err := existsCmd.Result()
	if err != nil {
		log.Errorf("exists %s err %+v", key, err)
		return false, err
	}
	err = hSetCmd.Err()
	if err != nil {
		log.Errorf("hset %s %+v err %+v", key, values, err)
		return false, err
	}
	err = expireCmd.Err()
	if err != nil {
		log.Warnf("expire %s err %+v", key, err)
	}

	return exists > 0, nil
}

type setResult struct {
	id     uint32
	exists *redis.IntCmd
	hset   *redis.IntCmd
	expire *redis.BoolCmd
}

func (r *Redis) BatchSetTentacle(ctx context.Context, infoList map[uint32]*pb.TentacleInfo) (map[uint32]bool, error) {
	resultList := make([]*setResult, 0, len(infoList))

	pipe := r.client.Pipeline()
	defer pipe.Close()

	for id, info := range infoList {
		key := r.genUserKey(id)

		values := make([]interface{}, 0, len(info.TentacleList)*2)
		for i := range info.TentacleList {
			values = append(values, info.TentacleList[i].Key, info.TentacleList[i].Val)
		}

		resultList = append(resultList, &setResult{
			id:     id,
			exists: pipe.Exists(ctx, key),
			hset:   pipe.HSet(ctx, key, values...),
			expire: pipe.Expire(ctx, key, r.expire),
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return nil, err
	}

	result := make(map[uint32]bool, len(infoList))
	for i := range resultList {
		exists, err := resultList[i].exists.Result()
		if err != nil {
			log.Warnf("exists [%d] err %+v", resultList[i].id, err)
			continue
		}
		err = resultList[i].hset.Err()
		if err != nil {
			log.Warnf("hset [%d] %+v err %+v", resultList[i].id, infoList[resultList[i].id], err)
			continue
		}
		err = resultList[i].expire.Err()
		if err != nil {
			log.Warnf("expire [%d] err %+v", resultList[i].id, err)
			continue
		}
		result[resultList[i].id] = exists > 0
	}
	return result, nil
}

func (r *Redis) GetTentacle(ctx context.Context, id uint32) (map[string]string, error) {
	key := r.genUserKey(id)
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
		result = append(result, pipe.HGetAll(ctx, r.genUserKey(idList[i])))
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
			log.Warnf("hgetall %s err %+v", r.genUserKey(idList[i]), err)
			continue
		}
		resultList[idList[i]] = result
	}
	return resultList, nil
}

func (r *Redis) DelTentacle(ctx context.Context, id uint32, infoList []*pb.Tentacle) (bool, error) {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	key := r.genUserKey(id)
	fields := make([]string, 0, len(infoList))
	for i := range infoList {
		fields = append(fields, infoList[i].Key)
	}
	hDelCmd := r.client.HDel(ctx, key, fields...)
	existsCmd := r.client.Exists(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return false, err
	}
	err = hDelCmd.Err()
	if err != nil {
		log.Errorf("hdel %s %+v err %+v", key, fields, err)
		return false, err
	}
	exists, err := existsCmd.Result()
	if err != nil {
		log.Errorf("exists %s err %+v", key, err)
		return false, err
	}

	return exists > 0, nil
}

type delResult struct {
	id     uint32
	del    *redis.IntCmd
	exists *redis.IntCmd
}

func (r *Redis) BatchDelTentacle(ctx context.Context, infoList map[uint32]*pb.TentacleInfo) (map[uint32]bool, error) {
	resultList := make([]*delResult, 0, len(infoList))

	pipe := r.client.Pipeline()
	defer pipe.Close()

	for id, info := range infoList {
		key := r.genUserKey(id)

		values := make([]string, 0, len(info.TentacleList))
		for i := range info.TentacleList {
			values = append(values, info.TentacleList[i].Key)
		}

		resultList = append(resultList, &delResult{
			id:     id,
			del:    pipe.HDel(ctx, key, values...),
			exists: pipe.Exists(ctx, key),
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return nil, err
	}

	result := make(map[uint32]bool, len(infoList))
	for i := range resultList {
		err := resultList[i].del.Err()
		if err != nil {
			log.Warnf("hel [%d] %+v err %+v", resultList[i].id, infoList[resultList[i].id], err)
			continue
		}
		exists, err := resultList[i].exists.Result()
		if err != nil {
			log.Errorf("exists [%d] err %+v", resultList[i].id, err)
			continue
		}
		result[resultList[i].id] = exists > 0
	}
	return result, nil
}

func (r *Redis) SetTentacleHistory(ctx context.Context, id uint32, record string, ts time.Time) error {
	pipe := r.client.Pipeline()
	defer pipe.Close()

	key := r.genHistoryKey(id)
	z := &redis.Z{
		Score:  float64(ts.UnixMilli()),
		Member: record,
	}
	zaddcmd := pipe.ZAdd(ctx, key, z)
	intCmd := pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10))
	expireCmd := pipe.Expire(ctx, key, time.Hour*24*7)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("exec err %+v", err)
		return err
	}

	err = zaddcmd.Err()
	if err != nil {
		log.Errorf("zadd %s %+v err %+v", key, z, err)
		return err
	}
	err = intCmd.Err()
	if err != nil {
		log.Warnf("zremrangebyscore %s err %+v", key, err)
	}
	err = expireCmd.Err()
	if err != nil {
		log.Warnf("expire %s err %+v", key, err)
	}

	return nil
}

func (r *Redis) GetTentacleHistory(ctx context.Context, id uint32, from, to time.Time) ([]string, error) {
	key := r.genHistoryKey(id)
	result, err := r.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(from.UnixMilli(), 10),
		Max: strconv.FormatInt(to.UnixMilli(), 10),
	}).Result()

	if err != nil {
		log.Errorf("zrangebyscore %s err %+v", key, err)
		return nil, err
	}
	return result, nil
}

func (r *Redis) Close(ctx context.Context) error {
	return r.client.Close()
}

// ---------------------------------------------------------------------------------------------------------------------

func (r *Redis) genUserKey(id uint32) string {
	return fmt.Sprintf("%s@%d", r.prefix, id)
}

func (r *Redis) genHistoryKey(id uint32) string {
	return fmt.Sprintf("%s@history@%d", r.prefix, id)
}
