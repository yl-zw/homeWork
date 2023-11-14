package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"webbook/config"
	_ "webbook/config"
	"webbook/internal/respository"
	"webbook/internal/respository/dao"
	"webbook/internal/respository/dao/user"
	service2 "webbook/internal/service"
	"webbook/internal/web"
	"webbook/internal/web/middleware"
)

func main() {
	err := dao.InitTables()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := config.NewDb()
	if err != nil {
		fmt.Println(err)
		return
	}

	userDao := user.NewUserDao(db)
	userRepository := respository.NewUserRepository(userDao)
	useService := service2.NewUseService(userRepository)
	router := web.New(useService)
	service := gin.Default()
	middleware.Init(service)
	router.RegisterRoute(service)

	err = service.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
	}
}
