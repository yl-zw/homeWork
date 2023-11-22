package main

import (
	"github.com/gin-gonic/gin"
	"webbook/internal/respository"
	"webbook/internal/respository/cache"
	"webbook/internal/respository/dao/user"
	service2 "webbook/internal/service"
	"webbook/internal/web"
	"webbook/internal/web/middleware"
	"webbook/ioc"
)

func InitWebService() *gin.Engine {
	DB, client, err := ioc.InitConfig()
	if err != nil {
		panic(err)
		return nil
	}
	userDao := user.NewUserDao(DB)
	userCache := cache.NewUserCache(client)
	userRepository := respository.NewUserRepository(userDao, userCache)
	messageClient := ioc.NewThird()
	code := service2.NewCode(messageClient, userCache)
	useService := service2.NewUseService(userRepository, code)
	handel := web.New(useService)
	engine := gin.Default()
	middleware.InitMiddle()
	web.RegisterRoute(engine, handel)

	return engine
}
