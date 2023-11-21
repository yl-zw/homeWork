package service

import (
	"fmt"
	"math/rand"
	"webbook/ailiyunmsm/msm"
	"webbook/internal/respository/cache"
)

type CodeService struct {
	msmClient *msm.ALiMessageClient
	user      *cache.UserCache
}

func NewCode(client *msm.ALiMessageClient, user *cache.UserCache) *CodeService {
	return &CodeService{
		msmClient: client,
		user:      user,
	}
}

func (c *CodeService) SendCode(biz string, data interface{}, phone ...string) error {
	err := c.msmClient.Send("zw", code(), phone...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, v := range phone {
		err = c.user.Set(key(biz, v), data)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
func key(biz string, phone string) string {
	return fmt.Sprintf("user-%s-%s", biz, phone)
}
func code() string {
	intn := rand.Intn(100000)
	return fmt.Sprintf("%06d", intn)
}
func (c *CodeService) GetCode(biz, number string) (error, interface{}) {
	res, err := c.user.Get(key(biz, number))
	if err != nil {
		fmt.Println(err)
		return err, nil
	}
	return nil, res
}
