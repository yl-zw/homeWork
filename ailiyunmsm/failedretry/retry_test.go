package failedretry

import (
	"errors"
	"fmt"
	cache2 "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	codemock "webbook/ailiyunmsm/mocks"
	"webbook/internal/respository/cache"
	cachemock "webbook/internal/respository/cache/mocks"

	"webbook/limiter"
)

func TestMSMRetry_Send(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(controller *gomock.Controller) *MSMRetry
		after      func(controller *MSMRetry, key string)
		singerName string
		code       string
		phone      string
		wanterr    error
		want       interface{}
	}{
		{
			name: "同步发送成功",
			mock: func(controller *gomock.Controller) *MSMRetry {
				code := codemock.NewMockCode(controller)
				code.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				limit := limiter.NewLimit(1, 10)
				cache := cachemock.NewMockCache(controller)
				return NewMSMRetry(code, limit, cache, "127.0.0.1", 3)
			},
			singerName: "xxxx",
			code:       "123",
			phone:      "12344445",
			wanterr:    nil,
		},
		{
			name: "已限流通过异步发送",
			mock: func(controller *gomock.Controller) *MSMRetry {
				code := codemock.NewMockCode(controller)
				code.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(UnknownError).AnyTimes()
				limit := limiter.NewLimit(1, 0)

				c := cache2.New(time.Minute*10, time.Minute*10)
				localCache := cache.NewLocalCache(c)
				return NewMSMRetry(code, limit, localCache, "127.0.0.1", 3)
			},
			after: func(retry *MSMRetry, key string) {
				get, err := retry.localCache.Get(key)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(get)
			},
			singerName: "xxxx",
			code:       "123",
			phone:      "12344445",
			wanterr:    nil,
		},
		{
			name: "第三方服务器错误已转为异步发送",
			mock: func(controller *gomock.Controller) *MSMRetry {
				code := codemock.NewMockCode(controller)
				code.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(UnknownError).AnyTimes()
				limit := limiter.NewLimit(1, 50)
				c := cache2.New(time.Minute*10, time.Minute*10)
				localCache := cache.NewLocalCache(c)
				return NewMSMRetry(code, limit, localCache, "127.0.0.1", 3)
			},
			after: func(retry *MSMRetry, key string) {
				retry.localCache.DelKey(key)
			},
			singerName: "xxxx",
			code:       "123",
			phone:      "12344445",
			wanterr:    nil,
		},
		{
			name: "同步发送失败",
			mock: func(controller *gomock.Controller) *MSMRetry {
				code := codemock.NewMockCode(controller)
				code.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("sssss"))
				limit := limiter.NewLimit(1, 50)
				cache := cachemock.NewMockCache(controller)
				return NewMSMRetry(code, limit, cache, "127.0.0.1", 3)
			},
			singerName: "xxxx",
			code:       "123",
			phone:      "12344445",
			wanterr:    errors.New("发送失败"),
		},
	}

	for _, v := range testCase {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			//err := v.mock(ctrl).Send(v.singerName, v.phone, v.code)
			//assert.Error(t, err, v.wanterr)

			can := 1
			tt := v.mock(ctrl)

			for can >= 0 {

				err := tt.Send(v.singerName, v.code, v.phone)
				//assert.Error(t, err, v.wanterr)
				assert.NoError(t, err, v.wanterr)
				can--
			}
			//v.after(tt, v.phone)

		})
	}
}
