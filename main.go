package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"webbook/ailiyunmsm/msm"
	"webbook/config"
	_ "webbook/config"
	"webbook/internal/respository"
	"webbook/internal/respository/cache"
	"webbook/internal/respository/dao"
	"webbook/internal/respository/dao/user"
	service2 "webbook/internal/service"
	"webbook/internal/web"
	"webbook/internal/web/middleware"
	"webbook/limiter"
)

func main() {
	limit := limiter.NewLimit(1, 50)
	middleware.Limiters = limit
	err := dao.InitTables()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, red, err := config.NewDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	userDao := user.NewUserDao(db)
	userCache := cache.NewUserCache(red)
	userRepository := respository.NewUserRepository(userDao, userCache)
	clientMSM, err := msm.NewClientMSM("", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	code := service2.NewCode(clientMSM, userCache)

	useService := service2.NewUseService(userRepository, code)
	router := web.New(useService)
	service := gin.Default()
	middleware.Init(service)
	router.RegisterRoute(service)

	err = service.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
	}
}
