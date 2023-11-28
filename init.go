package main

import (
	"github.com/gin-gonic/gin"
	localCache "github.com/patrickmn/go-cache"
	"webbook/config"
	"webbook/internal/respository"
	"webbook/internal/respository/cache"
	"webbook/internal/respository/dao/user"
	service2 "webbook/internal/service"
	"webbook/internal/web"
	"webbook/internal/web/middleware"
	"webbook/ioc"
	"webbook/logger"
)

var save = make(map[string]localCache.Item)

func InitWebService() *gin.Engine {
	DB, client, err := ioc.InitConfig()
	if err != nil {
		panic(err)
		return nil
	}
	logger := logger.InitLogger(config.Con.DB.Logger.IsDebug)
	userDao := user.NewUserDao(DB)
	userCache := cache.NewRedsCache(client)
	//localCache := localCache.NewFrom(localCache.DefaultExpiration, time.Minute*20, save)
	//userCache := cache.NewLocalCache(localCache)
	userRepository := respository.NewUserRepository(userDao, userCache)
	messageClient := ioc.NewThird()
	code := service2.NewCode(messageClient, userCache)
	useService := service2.NewUseService(userRepository, code, logger)
	handel := web.New(useService)
	engine := gin.Default()
	middleware.InitMiddle(engine, logger)
	web.RegisterRoute(engine, handel)

	return engine
}
