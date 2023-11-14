package web

import (
	"github.com/gin-gonic/gin"
	"webbook/internal/service"
	"webbook/internal/web/middleware"
)

type ServiceHandel struct {
	UserService *service.UseService
}

func New(useService *service.UseService) *ServiceHandel {
	return &ServiceHandel{UserService: useService}
}
func (S *ServiceHandel) RegisterRoute(service *gin.Engine) {
	// 使用中间件解决跨域问题
	service.Use(middleware.Cors(), middleware.CheckLogin())

	{
		user := service.Group("/user")
		user.POST("/login", S.UserService.Login)
		user.POST("/signup", S.UserService.SignUp)
		user.POST("/edit", S.UserService.Edit)
		user.GET("/profile", S.UserService.Profile)
	}

}
