package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"webbook/internal/domain"
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

func (u *UserCache) Set(uid string, user domain.Profile) error {
	res, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return u.cmd.Set(u.key(uid), res, u.Expiration).Err()
}

func (u *UserCache) Get(uid string) (domain.Profile, error) {
	key := u.key(uid)
	result, err := u.cmd.Get(key).Result()
	if err != nil {
		return domain.Profile{}, err
	}
	res := domain.Profile{}
	err = json.Unmarshal([]byte(result), &res)
	if err != nil {
		return domain.Profile{}, err
	}
	return res, nil

}
func (u *UserCache) key(uid string) string {
	return fmt.Sprintf("user-info-%s", uid)
}
