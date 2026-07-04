package sinking_web

import (
	"sync"
	"time"
)

type LimitRate struct {
	rate            int
	limit           int
	currentAmount   int
	lastConsumeTime int64
	lock            sync.Mutex
}

func (*LimitRate) currentTime() int64 {
	return time.Now().Unix()
}

func (tb *LimitRate) Wait(n int) {
	if n > tb.limit {
		return
	}
	for {
		tb.lock.Lock()
		now := tb.currentTime()
		if now == tb.lastConsumeTime {
			pre := tb.currentAmount + int(now-tb.lastConsumeTime)*tb.rate
			if pre > tb.limit {
				tb.currentAmount = tb.limit
			} else {
				tb.currentAmount = pre
			}
		} else {
			tb.currentAmount = tb.limit
		}
		if n <= tb.currentAmount {
			tb.currentAmount -= n
			tb.lastConsumeTime = now
			tb.lock.Unlock()
			return
		}
		tb.lock.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func (tb *LimitRate) Check(n int) bool {
	if n > tb.limit {
		return false
	}
	tb.lock.Lock()
	defer tb.lock.Unlock()
	now := tb.currentTime()
	res := false
	if now == tb.lastConsumeTime {
		if tb.currentAmount <= 0 {
			res = true
		}
	} else {
		tb.currentAmount = tb.limit
	}
	tb.currentAmount -= n
	tb.lastConsumeTime = now
	return res
}

var (
	limitRates     = make(map[string]*LimitRate)
	limitRatesLock sync.Mutex
)

func GetLimitRateIns(key string, limit int) *LimitRate {
	limitRatesLock.Lock()
	defer limitRatesLock.Unlock()
	obj := limitRates[key]
	if obj == nil {
		obj = NewLimitRate(limit, limit)
		limitRates[key] = obj
	}
	return obj
}

func NewLimitRate(limit int, rate int) *LimitRate {
	tb := LimitRate{
		rate:          rate,
		limit:         limit,
		currentAmount: limit,
	}
	tb.lastConsumeTime = tb.currentTime()
	return &tb
}
