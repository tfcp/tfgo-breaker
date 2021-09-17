# tfgo-breaker

## 1. Intro
This is a easy breaker by golang code. U can use it in your project quickly.
Support function break, timeout, auto dry-run.

## 2. Demo

```
    func main() {
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
```

## 3. Notice
Our breaker should block a process or a function

## 4. breaker conf
|  title   | describe  |default|
|  ----  | ----  |----|
|  Name | breaker name. one business process or function have one breaker.    |  |
|  Threshold | breaker reached value. key include (jsonRpcUrl+jsonRpc)  | 500 |
|  Expire  | breakerKey expired time. |5 min|
|  DryRunPercent  | when breaker is opened, there is ?% requests can pass breaker, if one of those success, breaker will closed. (example: 10: 1/10  5: 1/5 2: 1/2) |10%|
|  CallBackFunc  | normal business logic. when breaker is closed, will run it | |
|  BreakerFunc   | breaker business logic. when breaker is opened, will run it| |
