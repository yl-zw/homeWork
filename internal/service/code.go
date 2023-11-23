package service

import (
	"fmt"
	"math/rand"
	"webbook/ailiyunmsm"
	"webbook/internal/respository/cache"
)

type CodeS interface {
	SendCode(biz string, data interface{}, phone string) error
	GetCode(biz, number string, other ...string) (error, interface{})
	DelKeyByName(biz, phone string, other ...string) error
}

type CodeService struct {
	msmClient ailiyunmsm.Code
	user      cache.Cache
}

func NewCode(client ailiyunmsm.Code, user cache.Cache) *CodeService {
	return &CodeService{
		msmClient: client,
		user:      user,
	}
}

func (c *CodeService) SendCode(biz string, data interface{}, phone string) error {
	codes := code()
	fmt.Println(codes)
	//err := c.msmClient.Send("张巍的博客", codes, phone)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	fmt.Println(codes)
	err := c.user.Set(key(biz, phone), codes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(key(biz, phone))
	err = c.user.Set(key(biz, phone, "info"), data)
	if err != nil {
		fmt.Println(err)
		fmt.Println(key(biz, phone, "info"))
	}
	return nil
}
func key(biz string, phone string, other ...string) string {
	return fmt.Sprintf("user-%s-%s%s", biz, phone, other)
}
func code() string {
	intn := rand.Intn(99999)
	ss := fmt.Sprintf("%05d", intn)
	fmt.Println(ss)
	return fmt.Sprintf("%05d", intn)
}
func (c *CodeService) GetCode(biz, number string, other ...string) (error, interface{}) {
	res, err := c.user.Get(key(biz, number, other...))
	if err != nil {
		fmt.Println(err)
		return err, nil
	}
	return nil, res
}
func (c *CodeService) DelKeyByName(biz, phone string, other ...string) error {
	s := key(biz, phone, other...)
	fmt.Println(s)
	err := c.user.DelKey(s)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
