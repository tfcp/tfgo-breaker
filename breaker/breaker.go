package breaker

import (
	errs "github.com/pkg/errors"
	"time"
)

type Breaker struct {
	conf *BreakConf
}

func NewBreaker(conf *BreakConf) *Breaker {
	return &Breaker{conf: conf}
}

// get breaker status (true: breaker opened  false: breaker closed)
func (this *Breaker) GetBreakerStatus() bool {
	breakerKey := getCacheKey(this.conf.Name)
	breakValue, _ := cache.Get(breakerKey)
	if breakValue == nil {
		breakValue = 0
	}
	if breakValue == this.conf.Threshold {
		return true
	}
	return false
}

// do breaker logic
func (this *Breaker) Run(request interface{}) (interface{}, error) {
	breakerKey := getCacheKey(this.conf.Name)
	breakValue, _ := cache.Get(breakerKey)
	if breakValue == nil {
		breakValue = 0
	}
	if breakValue == this.conf.Threshold {
		// breaker opened
		if !canWeDryRun(this.conf.DryRunPercent) {
			// do breaker logic
			return this.breakFunc(request)
		}
		// dry run (default 1% request)
		res, err, breakerRes := this.conf.CallBackFunc(request)
		if true == breakerRes {
			// dry run success we can close breaker
			cache.Set(breakerKey, 0, this.conf.Expire)
		}
		return res, err
	}
	res, err, breakerRes := this.conf.CallBackFunc(request)
	if !breakerRes {
		// if normal logic error breaker add 1
		cache.Set(breakerKey, (breakValue).(int)+1, this.conf.Expire)
	}
	return res, err
}

func (this *Breaker) callBackFunc(req interface{}) (interface{}, error, bool) {
	tm := make(chan uint32)
	var (
		res     interface{}
		err     error
		breaker bool
	)
	go func() {
		res, err, breaker = this.conf.CallBackFunc(req)
		tm <- 1
	}()
	select {
	case <-tm:
		return res, err, breaker
	case <-time.After(breakerFuncTimeout):
		errMsg := "callBackFunc is timeout..."
		return res, errs.New(errMsg), breaker
	}
}

func (this *Breaker) breakFunc(req interface{}) (interface{}, error) {
	tm := make(chan uint32)
	var (
		res interface{}
		err error
	)
	go func() {
		res, err = this.conf.BreakerFunc(req)
		tm <- 1
	}()
	select {
	case <-tm:
		return res, err
	case <-time.After(breakerFuncTimeout):
		errMsg := "breakFunc is timeout..."
		return res, errs.New(errMsg)
	}
}
