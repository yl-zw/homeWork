package main

import (
	"fmt"
	_ "webbook/config"
)

func main() {

	//err := dao.InitTables()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//db, red, err := config.NewDb()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//userDao := user.NewUserDao(db)
	//userCache := cache.NewUserCache(red)
	//userRepository := respository.NewUserRepository(userDao, userCache)
	//clientMSM, err := msm.NewClientMSM("", "")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//code := service2.NewCode(clientMSM, userCache)
	//
	//useService := service2.NewUseService(userRepository, code)
	//router := web.New(useService)
	//service := gin.Default()
	service := InitWebService()
	//router.RegisterRoute(service)

	err := service.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
	}
}
