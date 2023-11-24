package failedretry

import (
	"errors"
	"fmt"
	"time"
	"webbook/ailiyunmsm"
	"webbook/internal/respository/cache"
	"webbook/limiter"
)

var (
	InternalError = errors.New("服务器内部错误")
	Throttling    = errors.New("时间段类限流")
	UnknownError  = errors.New("未知错误")
)

type MSMRetry struct {
	msmClient  ailiyunmsm.Code
	limit      *limiter.Limiter
	ip         string
	retryInfo  chan info
	canRetry   int
	localCache cache.Cache
}

type info struct {
	singerName  string
	code        string
	phoneNumber string
}

func NewMSMRetry(Client ailiyunmsm.Code, limiter2 *limiter.Limiter, localCache cache.Cache, ip string, reCan int) *MSMRetry {
	client := &MSMRetry{
		msmClient:  Client,
		limit:      limiter2,
		ip:         ip,
		retryInfo:  make(chan info, 10),
		canRetry:   reCan,
		localCache: localCache,
	}
	go func() {
		client.asySend()
	}()
	return client
}

func (M *MSMRetry) Send(singerName, code string, phoneNumber ...string) error {
	isLimit := M.limit.IsLimit(M.ip)
	if isLimit {
		go func() {
			M.retryInfo <- info{
				singerName:  singerName,
				code:        code,
				phoneNumber: phoneNumber[0],
			}
		}()
		fmt.Println("已限流异步发送")
		return nil
	}
	err := M.msmClient.Send(singerName, code, phoneNumber...)
	switch err {
	case nil:
		return err
	case InternalError, UnknownError, Throttling:
		M.retryInfo <- info{
			singerName:  singerName,
			code:        code,
			phoneNumber: phoneNumber[0],
		}
		return nil
	default:
		return errors.New("发送失败")

	}

}

func (M *MSMRetry) asySend() {
	var keys []string
	for true {
		//time.Sleep(time.Second)
		if infos, ok := <-M.retryInfo; ok {
			for M.canRetry != 0 {
				M.canRetry--
				err := M.msmClient.Send(infos.singerName, infos.code, infos.phoneNumber)
				if err != nil {
					fmt.Println(err)
				}
			}
			err := M.localCache.Set(infos.phoneNumber, infos, time.Minute*5)
			if err != nil {
				M.retryInfo <- infos
				fmt.Println(err)
				continue
			}
			keys = append(keys, infos.phoneNumber)
			continue

		}
		for k, v := range keys {
			reslut, err := M.localCache.Get(v)
			if err != nil {
				fmt.Println(err)
				continue
			}
			infoss := reslut.(info)
			err = M.msmClient.Send(infoss.singerName, infoss.code, infoss.phoneNumber)
			if err == nil {
				err = M.localCache.DelKey(v)
				if err != nil {
					fmt.Println(err)
					keys[k] = "nil"
				}
			}
		}

	}
}
