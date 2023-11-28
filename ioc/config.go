package ioc

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"webbook/config"
)

func InitConfig() (*gorm.DB, *redis.Client, error) {
	config.InitConfig()
	err := config.InitTables()
	if err != nil {
		panic(err)
		return nil, nil, err
	}
	DB, RedisClient, err := config.NewDb()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return DB, RedisClient, nil
}
