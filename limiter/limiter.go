package limiter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"time"
)

type Limiter struct {
	Pool map[string]info
	lock *sync.Mutex
	rate int64
	num  int64
}
type info struct {
	num  int64
	last int64
	max  int64
}

func NewLimit(rat int64, num int64) *Limiter {
	lim := &Limiter{
		Pool: make(map[string]info),
		lock: &sync.Mutex{},
		rate: rat,
		num:  num,
	}
	go func() {
		for {
			time.Sleep(time.Duration(lim.rate) * time.Second * 30)
			lim.flash()
			//fmt.Println(len(lim.Pool))

		}
	}()
	return lim
}

func (l *Limiter) IsLimit(ctx *gin.Context, now int64) bool {
	ip := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	l.lock.Lock()
	defer l.lock.Unlock()
	if val, ok := l.Pool[ip]; ok {
		//fmt.Println(strings.Split(ctx.Request.RemoteAddr, ":")[0])
		if val.num < l.num && now <= val.last {
			//fmt.Println(l.num)
			val.num++
			fmt.Println(val.num)
			l.Pool[ip] = val
			return false
		} else {
			if now > val.last {
				fmt.Println("刷新")
				l.Pool[ip] = info{
					last: now + l.rate*1000000000,
				}
				return false
			}
			return true
		}
	} else {
		l.Pool[ip] = info{
			last: now + l.rate*1000000000,
			num:  1,
		}
	}
	return false
}
func (l *Limiter) flash() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(l.Pool) < 0 {
		return
	}
	for k, v := range l.Pool {
		if (time.Now().UnixNano()-v.last)/1000000 > l.rate*30 {
			delete(l.Pool, k)
		}
	}
	return
}
