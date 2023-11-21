package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	http2 "webbook/http"
	"webbook/limiter"
)

const SessionIDKeyName = "userId"

var Limiters *limiter.Limiter

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cors.New(cors.Config{
			AllowedHeaders:   []string{"Content-type", "Auth", "User-Agent"},
			ExposedHeaders:   []string{"jwt-token"},
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
		if url.Path == "/user/signup" || url.Path == "/user/login" || url.Path == "/user/send" {
			return
		}
		//sessionCheck(ctx) //通过session身份验证
		tokenCheck(ctx) //通过token身份验证

	}
}

func LimiterMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//fmt.Println(len(Limiters.Pool))
		if Limiters.IsLimit(ctx, time.Now().UnixNano()) {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}
func Init(engine *gin.Engine) {
	//initSession(engine)
}

const updateTime = "updateTime"

func initSession(engine *gin.Engine) {
	store, err := sessions.NewRedisStore(10, "tcp", "192.168.3.108:6379", "", []byte("weefweif"))
	if err != nil {
		fmt.Println(err)
		return
	}
	engine.Use(sessions.Sessions("ssid", store))
}
func sessionCheck(ctx *gin.Context) {
	gob.Register(time.Now())
	val := sessions.Default(ctx).Get(SessionIDKeyName)
	if val == nil {
		var resp http2.Response
		resp.Code = http.StatusNonAuthoritativeInfo
		resp.Msg = "用户未登录"
		resp.Responses(ctx)
		return
	}
	if sessions.Default(ctx).Get(updateTime) == nil {
		sessions.Default(ctx).Set(updateTime, time.Now())
		sessions.Default(ctx).Set(SessionIDKeyName, val)
		sessions.Default(ctx).Options(sessions.Options{MaxAge: 900})
		err := sessions.Default(ctx).Save()
		if err != nil {
			fmt.Println(err)
		}
	}
	if time.Now().Sub(sessions.Default(ctx).Get(updateTime).(time.Time)) > time.Second*10 {
		sessions.Default(ctx).Set(updateTime, time.Now())
		sessions.Default(ctx).Set(SessionIDKeyName, val)
		sessions.Default(ctx).Options(sessions.Options{MaxAge: 900})
		err := sessions.Default(ctx).Save()
		if err != nil {
			fmt.Println(err)
		}

	}

}
func tokenCheck(ctx *gin.Context) {
	code := ctx.GetHeader("Auth")
	ug := UserClaim{}
	token, err := jwt.ParseWithClaims(code, &ug, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTkey), nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if !token.Valid {
		fmt.Println("身份过期")
		return
	}
	ctx.Set("email", ug.UserEmail)
	expirationTime, err := ug.GetExpirationTime()
	if err != nil {
		fmt.Println(err)
		return
	}
	if ug.UserAgent != ug.UserAgent {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if expirationTime.Sub(time.Now()) > time.Second*10 {
		ug.ExpiresAt = jwt.NewNumericDate(time.Now())
		newtoken, err := token.SignedString([]byte(JWTkey))
		if err != nil {
			fmt.Println(err)
		}
		ctx.Header("jwt-token", newtoken)
	}
}

type UserClaim struct {
	jwt.RegisteredClaims
	UserEmail string
	UserAgent string
}

var JWTkey = `vGUssghAQBrc6B887vPDQfBjdNhe2hh4`
