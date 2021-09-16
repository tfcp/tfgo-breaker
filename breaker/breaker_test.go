package breaker

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestNewBreaker(t *testing.T) {
	myNormalLogic := normalLogic
	myBreakerLogic := breakerLogic
	// create a new breaker
	breakerConf1 := NewBreakConf("test-breaker-1", 2, 20*time.Second, 2, myNormalLogic, myBreakerLogic)
	breaker1 := NewBreaker(breakerConf1)
	for i := 0; i < 100; i++ {
		r, _ := rand.Int(rand.Reader, big.NewInt(2))
		res, err := breaker1.Run(r.Int64())
		fmt.Println("breakerRun:", res, err)
		time.Sleep(1 * time.Second)
	}
}

func TestMultiBreaker(t *testing.T) {
	myNormalLogic := normalLogic
	myBreakerLogic := breakerLogic
	// create a new breaker
	breakerConf1 := NewBreakConf("test-breaker-1", 2, 20*time.Second, 2, myNormalLogic, myBreakerLogic)
	breaker1 := NewBreaker(breakerConf1)
	// create a new breaker
	breakerConf2 := NewBreakConf("test-breaker-2", 5, 20*time.Second, 2, myNormalLogic, myBreakerLogic)
	breaker2 := NewBreaker(breakerConf2)
	for i := 0; i < 100; i++ {
		r1, _ := rand.Int(rand.Reader, big.NewInt(2))
		res1, err1 := breaker1.Run(r1.Int64())
		fmt.Println("breakerRun1:", res1, err1)
		r2, _ := rand.Int(rand.Reader, big.NewInt(3))
		res2, err2 := breaker2.Run(r2.Int64())
		fmt.Println("breakerRun2:", res2, err2)
		time.Sleep(1 * time.Second)
	}
}

func normalLogic(request interface{}) (interface{}, error, bool) {
	// our normal logic
	if request.(int64) == 1 {
		errMsg := "normal logic err"
		return request, errors.New(errMsg), false
	}
	return "normal logic success", nil, true
}

func breakerLogic(request interface{}) (interface{}, error) {
	// our breaker logic
	return "this is breaker logic", nil
}
