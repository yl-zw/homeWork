package web

import (
	"github.com/gin-gonic/gin"
	"webbook/internal/service"
	"webbook/internal/web/middleware"
)

type ServiceHandel struct {
	UserService service.UserSer
}

func New(useService service.UserSer) *ServiceHandel {
	return &ServiceHandel{UserService: useService}
}
func RegisterRoute(service *gin.Engine, S *ServiceHandel) {
	// 使用中间件解决跨域问题
	service.Use(middleware.LimiterMiddle(), middleware.Cors(), middleware.CheckLogin(), middleware.LoggerMiddle())

	{
		user := service.Group("/user")
		user.POST("/login", S.UserService.Login)
		user.POST("/signup", S.UserService.SignUp)
		user.POST("/edit", S.UserService.Edit)
		user.GET("/profile", S.UserService.Profile)
		user.POST("/send", S.UserService.Send)
	}

}
