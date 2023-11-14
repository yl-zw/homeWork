package middleware

import (
	"fmt"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const SessionIDKeyName = "userId"

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cors.New(cors.Config{
			AllowedHeaders:   []string{"Content-type"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		})
	}
}

func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.URL
		if url.Path == "/user/signup" || url.Path == "/user/login" {
			return
		}
		if sessions.Default(ctx).Get(SessionIDKeyName) == nil {
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		} else {
			fmt.Println(sessions.Default(ctx).Get(SessionIDKeyName))
		}

	}
}
func Init(engine *gin.Engine) {
	store := sessions.NewCookieStore([]byte("secret"))
	engine.Use(sessions.Sessions("ssid", store))
}
