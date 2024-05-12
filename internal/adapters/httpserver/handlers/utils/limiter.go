package utils

import (
	"myapp/internal/config"
	"myapp/internal/core/domain"
	"sync"
	"time"
)

type Limiter struct {
	cl *ConcurrencyLimiter
	rl *RateLimiter
}

type ConcurrencyLimiter struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

type RateLimiter struct {
	UserRequests map[uint64]UsersRequests
	Limit        int
	Interval     time.Duration
	Mutex        *sync.Mutex
}

type UsersRequests struct {
	CountRequests int
	LastRequest   time.Time
}

func NewLimiter(httpCfg *config.HttpServerConfig) *Limiter {
	return &Limiter{
		NewConcurrencyLimiter(httpCfg.ConcurrencyLimit),
		NewRateLimiter(httpCfg.RateLimit),
	}
}

func NewConcurrencyLimiter(maxConcurrent int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		sem: make(chan struct{}, maxConcurrent),
	}
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		UserRequests: make(map[uint64]UsersRequests),
		Limit:        limit,
		Interval:     time.Second,
		Mutex:        &sync.Mutex{},
	}
}

func (cl *ConcurrencyLimiter) Add() {
	cl.sem <- struct{}{}
	cl.wg.Add(1)
	return
}

func (cl *ConcurrencyLimiter) Done() {
	<-cl.sem
	cl.wg.Done()
}

func (rl *RateLimiter) Add(id uint64) error {

	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	now := time.Now()

	userReq, _ := rl.UserRequests[id]
	if now.Sub(userReq.LastRequest) >= rl.Interval {
		rl.UserRequests[id] = UsersRequests{
			CountRequests: 1,
			LastRequest:   now,
		}
		return nil
	}

	if userReq.CountRequests >= rl.Limit {
		return domain.ErrRateLimitExceeded
	}

	rl.UserRequests[id] = UsersRequests{
		CountRequests: userReq.CountRequests + 1,
		LastRequest:   now,
	}
	return nil
}
