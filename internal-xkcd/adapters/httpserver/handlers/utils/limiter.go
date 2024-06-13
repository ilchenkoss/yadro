package utils

import (
	"myapp/internal-xkcd/config"
	"myapp/internal-xkcd/core/domain"
	"sync"
	"time"
)

type Limiter struct {
	cl *ConcurrencyLimiter
	rl *RateLimiter
}

type ConcurrencyLimiter struct {
	sem chan struct{}
}

type RateLimiter struct {
	UserRequests map[int64]UsersRequests
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
		UserRequests: make(map[int64]UsersRequests),
		Limit:        limit,
		Interval:     time.Second,
		Mutex:        &sync.Mutex{},
	}
}

func (cl *ConcurrencyLimiter) Add() {
	cl.sem <- struct{}{}
}

func (cl *ConcurrencyLimiter) Done() {
	<-cl.sem
}

func (rl *RateLimiter) Add(id int64) error {

	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	now := time.Now()

	userReq := rl.UserRequests[id]
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
