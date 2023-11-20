package respository

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
	"webbook/internal/domain"
	"webbook/internal/respository/cache"
	"webbook/internal/respository/dao/user"
)

type UserRepository struct {
	dao   *user.UseDao
	cache *cache.UserCache
}

var (
	ErrUserEmailErr    = user.EmailError
	ErrEmailOrPassword = user.RecodNotFound
	RedisErr           = cache.KeyIsNotExist
)

func NewUserRepository(useDao *user.UseDao, userCache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   useDao,
		cache: userCache,
	}
}
func (u *UserRepository) Create(ctx context.Context, in *domain.User) error {
	generateFromPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.dao.Insert(ctx, user.User{
		Email:    in.Email,
		Password: string(generateFromPassword),
	})
}

func (u *UserRepository) GetUserInfo(ctx *gin.Context, req *domain.ReqLoginUser) (string, error) {
	var user = &domain.User{}
	uinfo, err := u.dao.GetUserInfoByEmail(ctx, req)
	if err != nil {
		return "", err
	}

	u.ToDomain(&uinfo, user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println(err)
		return "", ErrEmailOrPassword
	}

	return user.Email, nil
}

func (u *UserRepository) GetProfileInfo(ctx *gin.Context, email interface{}) (domain.Profile, error) {
	var res domain.Profile
	tem := email.(string)
	res, err := u.cache.Get(tem)
	if err != nil {
		if err != RedisErr {
			return res, err
		}
		fmt.Println(err)
	}
	info, err := u.dao.GetProfileInfo(ctx, tem)
	if err != nil {
		return res, err
	}

	res.Email = info.Email
	format := time.Unix(info.Birthday, 0).Format("2006-01-02")
	res.Birthday = format
	res.PersonalProfile = info.Personalprofile
	res.UserName = info.Username
	err = u.cache.Set(tem, res)
	if err != nil {
		fmt.Println(err)
	}
	return res, nil
}

func (u *UserRepository) CreateProfile(ctx *gin.Context, profile *domain.Profile) error {
	var res = &user.Profile{}
	res.Email = profile.Email
	res.Personalprofile = profile.PersonalProfile
	t, err := time.Parse("2006-01-02", profile.Birthday)
	if err != nil {
		return err
	}
	res.Birthday = t.Unix()
	return u.dao.InsertProfile(ctx, res)
}

func (u *UserRepository) UpdateProfile(ctx *gin.Context, req *domain.Profile) error {
	var updata user.Profile
	updata.Personalprofile = req.PersonalProfile
	updata.Username = req.UserName
	updata.Email = req.Email
	t, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return err
	}
	updata.Birthday = t.Unix()
	return u.dao.UpdateProfile(ctx, updata)
}

func (u *UserRepository) ToDomain(user *user.User, user2 *domain.User) {
	user2.ID = int(user.Id)
	user2.Email = user.Email
	user2.Password = user.Password
}
