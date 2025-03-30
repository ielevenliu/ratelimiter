package main

import (
	"sync"
	"time"
)

/*
* 固定窗口限流算法:
*   1.假设单位时间(windowUnit)是1s，限流阀值(threshold)为3
*   2.当次数少于限流阀值，就允许访问，并且计数器+1
*   3.当次数大于限流阀值，就拒绝访问
*   4.当前的时间窗口过去之后，计数器清零
* 存在 '临界问题'
*  临界问题: 假设限流阀值为3个请求,单位时间窗口是1s,
*    如果我们在单位时间内的前0.8-1s和1-1.2s,分别并发3个请求,
*    虽然都没有超过阀值,但是如果算0.8-1.2s则并发数高达10,已经超过单位时间1s不超过3的阀值了
 */
const (
	threshold  = 3
	windowUnit = 1
)

var (
	counter  int
	lastTime int
	once     sync.Once
)

func FixedWindowRateLimiter() bool {
	once.Do(func() {
		lastTime = time.Now().Second()
	})

	curTime := time.Now().Second()
	if curTime-lastTime > windowUnit {
		counter = 0
		lastTime = curTime
	}
	if counter < threshold {
		counter++
		return true
	}
	return false
}
