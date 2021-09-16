package breaker

import (
	"fmt"
	"github.com/gogf/gf/os/gcache"
	"math/rand"
	"time"
)

var (
	// default breaker threshold
	breakerThreshold = 500
	// breaker value expire time(default 10 min)
	breakerExpire = 10 * 60 * time.Second
	// breaker cacheKey
	breakerKey = "tfbreaker#%s"
	// dry-run percent (default 1% request can pass breaker)
	breakerDryRunPercent = 100
	// funcTimeOut
	breakerFuncTimeout = 5 * time.Second
	// global breaker cache(use gf cache component)
	cache = gcache.New()
)

type BreakConf struct {
	// breaker name (one function or process have one breaker)
	Name string
	// callbackFunc can return false max value. if breaker value over Threshold, we will open breaker and do BreakerFunc.
	// every request have own Threshold
	Threshold int
	//	breaker value expire-time. if request times is too low, cannot reach dry-run.
	Expire time.Duration
	// after breaker opened,we can allow ?% requests pass breaker,if dry-run success, close the breaker.
	DryRunPercent int
	// breaker callback
	// bool: need to mark breaker
	CallBackFunc func(request interface{}) (interface{}, error, bool)
	// breaker func
	BreakerFunc func(request interface{}) (interface{}, error)
}

// breaker conf
func NewBreakConf(breakName string, threshold int, expire time.Duration, dryRunPercent int,
	callFunc func(request interface{}) (interface{}, error, bool),
	breakerFunc func(request interface{}) (interface{}, error)) *BreakConf {
	breakerConf := &BreakConf{
		Name:          breakName,
		Threshold:     breakerThreshold,
		Expire:        breakerExpire,
		DryRunPercent: breakerDryRunPercent,
		CallBackFunc:  callFunc,
		BreakerFunc:   breakerFunc,
	}
	if threshold > 0 {
		breakerConf.Threshold = threshold
	}
	if expire >= 1*time.Second {
		breakerConf.Expire = expire
	}
	if dryRunPercent > 1 {
		breakerConf.DryRunPercent = dryRunPercent
	}
	return breakerConf
}

// 获取cacheKey
func getCacheKey(key string) string {
	res := fmt.Sprintf(breakerKey, key)
	return res
}

func canWeDryRun(dryRunPercent int) bool {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(dryRunPercent)
	if value == 1 {
		return true
	}
	return false
}
