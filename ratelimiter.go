package main

import "time"

/*
* 固定窗口限流算法:
*   1.假设单位时间(WindowInterval)是1s，限流阀值(Threshold)为3
*   2.当次数少于限流阀值，就允许访问，并且计数器+1
*   3.当次数大于限流阀值，就拒绝访问
*   4.当前的时间窗口过去之后，计数器清零
* 存在 '临界问题'
*  临界问题: 假设限流阀值为3个请求,单位时间窗口是1s,
*    如果我们在单位时间内的前0.8-1s和1-1.2s,分别并发3个请求,
*    虽然都没有超过阀值,但是如果算0.8-1.2s则并发数高达10,已经超过单位时间1s不超过3的阀值了
 */

type FixedWindowRateLimiter struct {
	Threshold      int // 限流阈值
	WindowInterval int // 窗口时间间隔
	Counter        int
	LastTime       int
}

func NewFixedWindowRateLimiter() *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		Threshold:      3,
		WindowInterval: 1,
		LastTime:       time.Now().Second(),
	}
}

func (f *FixedWindowRateLimiter) limiter() bool {
	curTime := time.Now().Second()
	if curTime-f.LastTime > f.WindowInterval {
		f.Counter = 0
		f.LastTime = curTime
	}
	if f.Counter < f.Threshold {
		f.Counter++
		return true
	}
	return false
}

/*
* 滑动窗口限流算法:
*  1.将单位时间周期分为n个小周期,分别记录每个小周期内接口的访问次数,并且根据时间滑动删除过期的小周期
* 优点:
*  1.滑动窗口算法解决了固定窗口的临界问题
*  2.当滑动窗口的格子周期划分的越多,那么滑动窗口的滚动就越平滑,限流的统计就会越精确
* 缺点:
*  1.一旦到达限流后,请求都会直接暴力被拒绝,我们会损失一部分请求,这其实对于产品来说,并不太友好
 */
type SlidingWindowRateLimiter struct {
	Threshold         int
	SubWindowNum      int64
	SubWindowInterval int64
	Counters          map[int64]int
}

func NewSlidingWindowRateLimiter() *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		Threshold:         3,
		SubWindowNum:      10,
		SubWindowInterval: 100,
		Counters:          make(map[int64]int),
	}
}

func (s *SlidingWindowRateLimiter) limiter() bool {
	// 获取当前小窗口
	curWindow := time.Now().UnixMilli() / s.SubWindowInterval * s.SubWindowInterval
	// 获取起始小窗口
	startWindow := curWindow - s.SubWindowInterval*(s.SubWindowNum-1)

	var totalCnt int
	for idx, cnt := range s.Counters {
		if idx < startWindow {
			delete(s.Counters, idx)
			continue
		}
		totalCnt += cnt
	}
	if totalCnt >= s.Threshold {
		return false
	}
	s.Counters[curWindow]++
	return true
}
