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
}

func (*LimitRate) currentTime() int64 {
	return time.Now().Unix()
}

func (tb *LimitRate) Wait(n int) {
	if n > tb.limit {
		return
	}
	if tb.currentTime() == tb.lastConsumeTime {
		ticker := time.NewTicker(500 * time.Millisecond)
		for n > tb.currentAmount {
			pre := tb.currentAmount + int(tb.currentTime()-tb.lastConsumeTime)*tb.rate
			if pre > tb.limit {
				tb.currentAmount = tb.limit
			} else {
				tb.currentAmount = pre
			}
			<-ticker.C
		}
	} else {
		tb.currentAmount = tb.limit
	}
	tb.currentAmount -= n
	tb.lastConsumeTime = tb.currentTime()
}

func (tb *LimitRate) Check(n int) bool {
	if n > tb.limit {
		return false
	}
	res := false
	if tb.currentTime() == tb.lastConsumeTime {
		if tb.currentAmount <= 0 {
			res = true
		}
	} else {
		tb.currentAmount = tb.limit
	}
	tb.currentAmount -= n
	tb.lastConsumeTime = tb.currentTime()
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
