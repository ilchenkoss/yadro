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

func (l *Limiter) Add(id uint64) error {

	l.rl.Mutex.Lock()
	defer l.rl.Mutex.Unlock()

	now := time.Now()

	userReq, _ := l.rl.UserRequests[id]
	if now.Sub(userReq.LastRequest) >= l.rl.Interval {
		l.rl.UserRequests[id] = UsersRequests{
			CountRequests: 1,
			LastRequest:   now,
		}
		l.cl.sem <- struct{}{}
		l.cl.wg.Add(1)
		return nil
	}

	if userReq.CountRequests >= l.rl.Limit {
		time.Sleep(time.Duration(5000) * time.Millisecond)
		return domain.ErrRateLimitExceeded
	}

	l.rl.UserRequests[id] = UsersRequests{
		CountRequests: userReq.CountRequests + 1,
		LastRequest:   now,
	}
	l.cl.sem <- struct{}{}
	l.cl.wg.Add(1)
	return nil
}

func (l *Limiter) Done() {
	<-l.cl.sem
	l.cl.wg.Done()
}
