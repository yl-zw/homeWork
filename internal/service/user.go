package service

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	http2 "webbook/http"
	"webbook/internal/domain"
	"webbook/internal/respository"
	"webbook/internal/web/middleware"
)

const (
	email    = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	password = `^(?=.*\d)(?=.*[A-z])[\da-zA-Z]{8,}$$`
	birthDay = `^d{4}-d{2}-d{2}$`
)

type UserSer interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Profile(ctx *gin.Context)
	Send(ctx *gin.Context)
}
type UseService struct {
	email       *regexp2.Regexp
	password    *regexp2.Regexp
	birthDay    *regexp2.Regexp
	repo        respository.Repository
	codeService CodeS
}

var (
	ErrEmail                 = respository.ErrUserEmailErr
	ErrEmailOrPaawordIsWrong = respository.ErrEmailOrPassword
	RedisErr                 = respository.RedisErr
)

func NewUseService(repository respository.Repository, code CodeS) *UseService {
	return &UseService{
		email:       regexp2.MustCompile(email, regexp2.None),
		password:    regexp2.MustCompile(password, regexp2.None),
		birthDay:    regexp2.MustCompile(birthDay, regexp2.None),
		repo:        repository,
		codeService: code,
	}
}
func (U *UseService) SignUp(ctx *gin.Context) {
	var reqUser = &domain.ReqSingUpUser{}
	var res = &http2.Response{}
	if ctx.Bind(reqUser) != nil {

		res.Code = http.StatusBadRequest
		res.Msg = errors.New("不是json格式").Error()
		res.Responses(ctx)
		return
	}
	isEmail, err := U.email.MatchString(reqUser.Email)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Msg = errors.New("系统错误").Error()
		res.Responses(ctx)
		return
	}
	if !isEmail {
		res.Code = http.StatusBadRequest
		res.Msg = errors.New("邮箱格式错误").Error()
		res.Responses(ctx)
		return
	}
	isPassword, err := U.password.MatchString(reqUser.Password)
	if err != nil {
		fmt.Println(err)
		res.Code = http.StatusBadRequest
		res.Msg = errors.New("系统错误").Error()
		res.Responses(ctx)
		return
	}
	if !isPassword {
		res.Code = http.StatusBadRequest
		res.Msg = errors.New("密码格式错误").Error()
		res.Responses(ctx)
		return
	}
	if reqUser.Password != reqUser.ConfirmPassword {
		res.Code = http.StatusBadRequest
		res.Msg = errors.New("两次密码不一致").Error()
		res.Responses(ctx)
		return
	}
	if err != nil {
		res.Code = http.StatusNotAcceptable
		res.Msg = errors.New("系统错误").Error()
		res.Responses(ctx)
	}
	var user = &domain.User{}
	user.Email = reqUser.Email
	user.Password = reqUser.Password
	user.Phone = reqUser.Phone
	err = U.codeService.SendCode("regist", user, user.Phone)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "系统错误"
		res.Responses(ctx)
		return
	}
	//err = U.repo.SetKey("regist", reqUser.Phone, 123)
	//if err != nil {
	//	fmt.Println(err)
	//	res.Code = http.StatusInternalServerError
	//	res.Msg = "系统错误"
	//	res.Responses(ctx)
	//	return
	//}
	//err = U.repo.SetKey("registinfo", reqUser.Phone, user)
	//if err != nil {
	//	fmt.Println(err)
	//	res.Code = http.StatusInternalServerError
	//	res.Msg = "系统错误"
	//	res.Responses(ctx)
	//	return
	//}
	res.Code = http.StatusOK
	res.Msg = "请输入验证码"
	res.Responses(ctx)
	return

	//err = U.repo.Create(ctx, user)
	//if errors.Is(err, ErrEmail) {
	//	res.Code = http.StatusBadRequest
	//	res.Msg = "账号已注册"
	//	res.Responses(ctx)
	//	return
	//}
	//if err != nil {
	//	res.Code = http.StatusNotExtended
	//	res.Msg = "系统错误,注册失败"
	//	res.Responses(ctx)
	//	return
	//}
	//res.Code = http.StatusOK
	//res.Msg = "注册成功"
	//res.Responses(ctx)
	//return
}
func (U *UseService) Login(ctx *gin.Context) {
	var req = &domain.ReqLoginUser{}
	var res = &http2.Response{}
	if ctx.Bind(req) != nil {
		res.Code = http.StatusBadRequest
		res.Msg = "数据格式不符合规范"
	}
	userEmail, err := U.repo.GetUserInfo(ctx, req)
	if err == ErrEmailOrPaawordIsWrong {
		res.Code = http.StatusFailedDependency
		res.Msg = "账号或者密码错误"
		res.Responses(ctx)
		return
	}
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Msg = "系统错误，登陆失败"
		res.Responses(ctx)
		return
	}
	//openSession(ctx, userEmail, res) //开启session
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.UserClaim{
		UserEmail: userEmail,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	signedString, err := token.SignedString([]byte(middleware.JWTkey))
	if err != nil {
		fmt.Println(err)
		res.Code = http.StatusInternalServerError
		res.Msg = "系统错误"
		res.Responses(ctx)
		return
	}
	ctx.Header("jwt-token", signedString)

	res.Code = http.StatusOK
	res.Msg = "登录成功"
	res.Responses(ctx)
	return
}
func (U *UseService) Edit(ctx *gin.Context) {
	type Edit struct {
		UserName      string `json:"userName"`
		Email         string `json:"email"`
		PersonProfile string `json:"personProfile"`
		BirthDay      string `json:"birthDay"`
	}
	var info Edit
	var respon = &http2.Response{}
	if ctx.Bind(&info) != nil {
		respon.Code = http.StatusBadRequest
		respon.Msg = "请求参数不符合规范"
		respon.Responses(ctx)
		return
	}
	if len(info.UserName) > 15 || len(info.UserName) <= 0 {
		respon.Code = http.StatusBadRequest
		respon.Msg = "昵称不符合规则"
		respon.Responses(ctx)
		return
	}
	if len(info.PersonProfile) > 200 {
		respon.Code = http.StatusBadRequest
		respon.Msg = "个人简介过长"
		respon.Responses(ctx)
		return
	}
	isBirthDay, err := U.birthDay.MatchString(info.BirthDay)
	if err != nil {
		respon.Code = http.StatusBadRequest
		respon.Msg = "系统错误"
		respon.Responses(ctx)
		return
	}
	if isBirthDay {
		respon.Code = http.StatusBadRequest
		respon.Msg = "生日日期格式不对"
		respon.Responses(ctx)
		return
	}
	//email := sessions.Default(ctx).Get(middleware.SessionIDKeyName)
	value, exists := ctx.Get("email")
	if !exists {
		fmt.Println("未拿到数据")
		return
	}
	var req domain.Profile
	req.Email = value.(string)
	req.Birthday = info.BirthDay
	req.UserName = info.UserName
	req.PersonalProfile = info.PersonProfile
	profileInfo, err := U.repo.GetProfileInfo(ctx, email)
	if err != nil {
		respon.Code = http.StatusBadRequest
		respon.Msg = "系统错误"
		respon.Responses(ctx)
		return
	}
	if profileInfo.Email != "" {
		err = U.repo.UpdateProfile(ctx, &req)
		if err != nil {
			respon.Code = http.StatusBadRequest
			respon.Msg = "编辑失败"
			respon.Responses(ctx)
			return
		}
		respon.Code = http.StatusOK
		respon.Msg = "成功"
		respon.Data = req
		respon.Responses(ctx)
		return
	}
	err = U.repo.CreateProfile(ctx, &req)
	if err != nil {
		respon.Code = http.StatusBadRequest
		respon.Msg = "编辑失败"
		respon.Responses(ctx)
		return
	}
	respon.Code = http.StatusOK
	respon.Msg = "编辑成功"
	respon.Data = req
	respon.Responses(ctx)
	return
}
func (U *UseService) Profile(ctx *gin.Context) {
	var resp = &http2.Response{}
	//session := sessions.Default(ctx)
	//email := session.Get(middleware.SessionIDKeyName)
	value, exists := ctx.Get("email")
	if !exists {
		fmt.Println("未拿到数据")
		return
	}
	info, err := U.repo.GetProfileInfo(ctx, value)
	if err != nil {
		fmt.Println(err)
		resp.Code = http.StatusBadRequest
		resp.Msg = "查询失败"
		resp.Responses(ctx)
		return
	}
	resp.Code = http.StatusOK
	resp.Msg = "成功"
	resp.Data = info
	resp.Responses(ctx)
	return

}

func (U *UseService) Send(ctx *gin.Context) {
	type req struct {
		Biz   string `json:"biz"`
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	respon := &http2.Response{}
	var reqs = req{}
	if ctx.Bind(&reqs) != nil {
		respon.Code = http.StatusBadRequest
		respon.Msg = "参数错误"
		respon.Responses(ctx)
		return
	}
	err, result := U.codeService.GetCode(reqs.Biz, reqs.Phone)
	if err != nil && err != RedisErr {
		respon.Code = http.StatusInternalServerError
		respon.Msg = "系统错误"
		respon.Responses(ctx)
		return
	}
	if result == nil {
		respon.Code = http.StatusBadRequest
		respon.Msg = "验证码过期"
		respon.Responses(ctx)
		return
	}

	if reqs.Code != result {
		respon.Code = http.StatusBadRequest
		respon.Msg = "验证码错误"
		respon.Responses(ctx)
		return
	}
	err = U.codeService.DelKeyByName(reqs.Biz, reqs.Phone)
	fmt.Println(err)
	err, info := U.codeService.GetCode(reqs.Biz, reqs.Phone, "info")
	if err != nil {
		fmt.Println(err)
		respon.Code = http.StatusInternalServerError
		respon.Msg = "系统错误"
		respon.Responses(ctx)
		return
	}
	err = U.codeService.DelKeyByName(reqs.Biz, reqs.Phone, "info")
	if err != nil {
		fmt.Println(err)
	}
	tt := info.(*domain.User)
	//rr := domain.User{}
	//if json.Unmarshal([]byte(tt), &rr) != nil {
	//	fmt.Println(json.Unmarshal([]byte(tt), &rr))
	//	return
	//}
	err = U.repo.Create(ctx, tt)
	if err != nil {
		fmt.Println(err)
		respon.Code = http.StatusInternalServerError
		respon.Msg = "系统错误"
		respon.Responses(ctx)
		return
	}
	respon.Code = http.StatusOK
	respon.Msg = "注册成功"
	respon.Responses(ctx)
	return
}

func openSession(ctx *gin.Context, val string, response *http2.Response) {
	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: 900,
	})
	session.Set(middleware.SessionIDKeyName, val)
	err := session.Save()
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Msg = "系统错误，登录失败"
		response.Responses(ctx)
		return
	}
}
