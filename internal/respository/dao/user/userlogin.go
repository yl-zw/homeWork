package user

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"webbook/internal/domain"
)

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}

func (u *User) GetTableName() string {
	return "users"
}

type UseDao struct {
	db *gorm.DB
}

var (
	EmailError    = errors.New("唯一主键重复")
	RecodNotFound = errors.New("数据不存在")
)

func NewUserDao(db *gorm.DB) *UseDao {
	db = db.Debug()
	return &UseDao{db: db}
}
func (u *UseDao) Insert(ctx context.Context, user User) error {
	err := u.db.Table(user.GetTableName()).Create(&user).Error
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			const uniqueIndexErrNo = 1062
			if me.Number == uniqueIndexErrNo {
				return EmailError
			}
		}
		return err
	}

	return nil
}

func (u *UseDao) GetUserInfoByEmail(ctx *gin.Context, req *domain.ReqLoginUser) (User, error) {
	var info User
	err := u.db.Model(&info).Where("email= ?", req.Email).First(&info).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, RecodNotFound
	}
	return info, err
}

func (u *UseDao) GetProfileInfo(ctx *gin.Context, email string) (Profile, error) {
	var profile Profile
	err := u.db.Table(profile.GetTableName()).Where("email=?", email).First(&profile).Error
	if err == gorm.ErrRecordNotFound {
		return profile, nil
	}
	return profile, nil
}
func (u *UseDao) InsertProfile(ctx context.Context, user *Profile) error {
	err := u.db.Table(user.GetTableName()).Create(user).Error
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			const uniqueIndexErrNo = 1062
			if me.Number == uniqueIndexErrNo {
				return EmailError
			}
		}
		return err
	}

	return nil
}

func (u *UseDao) UpdateProfile(ctx *gin.Context, updata Profile) error {
	return u.db.Table(updata.GetTableName()).Where("email=?", updata.Email).Updates(updata).Error
}
