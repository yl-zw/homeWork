package ioc

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"os"
	"webbook/config"
)

func InitConfig() (*gorm.DB, *redis.Client, error) {
	path := os.Args[0]
	if path == "" {
		fmt.Println("差环境变量:--path")
		return nil, nil, errors.New("差环境变量")
	}
	config.InitConfig(path)
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
