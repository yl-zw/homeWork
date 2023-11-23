package cache

import (
	"errors"
	"fmt"
	localCache "github.com/patrickmn/go-cache"
	"time"
)

type LocalCache struct {
	cache  *localCache.Cache
	expire time.Duration
}

func (l LocalCache) Set(key string, value interface{}, expire ...time.Duration) error {
	if len(expire) > 0 {
		l.expire = expire[0]
	}
	err := l.cache.Add(key, value, l.expire)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (l LocalCache) Get(key string) (interface{}, error) {
	res, isGetKey := l.cache.Get(key)
	if !isGetKey {
		fmt.Println("key不存在")
		return nil, errors.New("key不存在")
	}

	return res, nil
}

func (l LocalCache) DelKey(key string) error {
	l.cache.Delete(key)
	return nil
}

func NewLocalCache(cache *localCache.Cache) *LocalCache {
	return &LocalCache{
		cache:  cache,
		expire: time.Minute * 15,
	}
}
