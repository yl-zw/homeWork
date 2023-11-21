package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

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
	return u.cmd.Set(u.key(uid), res, u.Expiration).Err()
}

func (u *UserCache) Get(uid string) (interface{}, error) {
	key := u.key(uid)
	result, err := u.cmd.Get(key).Result()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return result, err
	}
	return result, nil

}
func (u *UserCache) key(uid string) string {
	return fmt.Sprintf("user-info-%s", uid)
}
