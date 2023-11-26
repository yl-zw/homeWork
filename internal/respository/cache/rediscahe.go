package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expire ...time.Duration) error
	Get(key string) ([]byte, error)
	DelKey(key string) error
}

type RedisCache struct {
	cmd        redis.Cmdable
	Expiration time.Duration
}

const KeyIsNotExist = redis.Nil

func NewRedsCache(cmd redis.Cmdable) *RedisCache {
	return &RedisCache{
		cmd:        cmd,
		Expiration: time.Minute * 15,
	}
}

func (u *RedisCache) Set(uid string, user interface{}, expire ...time.Duration) error {
	//res, err := json.Marshal(user)
	//fmt.Println(res)
	//if err != nil {
	//	return err
	//}
	fmt.Println(user)
	return u.cmd.Set(uid, user, u.Expiration).Err()
}

func (u *RedisCache) Get(uid string) ([]byte, error) {
	result, err := u.cmd.Get(uid).Bytes()
	if err != nil {
		return result, err
	}

	return result, nil

}
func (u *RedisCache) DelKey(key string) error {

	return u.cmd.Del(key).Err()
}
