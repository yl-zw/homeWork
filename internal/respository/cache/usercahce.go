package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type Cache interface {
	Set(uid string, user interface{}) error
	Get(uid string) (interface{}, error)
	DelKey(key string) error
}

type UserCache struct {
	cmd        redis.Cmdable
	Expiration time.Duration
}

const KeyIsNotExist = redis.Nil

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		Expiration: time.Minute * 15,
	}
}

func (u *UserCache) Set(uid string, user interface{}) error {
	res, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return u.cmd.Set(uid, res, u.Expiration).Err()
}

func (u *UserCache) Get(uid string) (interface{}, error) {
	fmt.Println(uid)
	result, err := u.cmd.Get(uid).Result()
	if err != nil {
		return result, err
	}
	unquote, err := strconv.Unquote(result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return unquote, nil

}
func (u *UserCache) DelKey(key string) error {

	return u.cmd.Del(key).Err()
}
