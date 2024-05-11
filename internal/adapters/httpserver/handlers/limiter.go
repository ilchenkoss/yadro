package handlers

import (
	"myapp/internal/config"
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
	UserRequests map[int]UsersRequests
	Limit        int
	Interval     time.Duration
	Mutex        sync.Mutex
}

type UsersRequests struct {
	CountRequests int
	LastRequest   time.Time
}

type RateLimitExceededError struct {
	Message string
}

func (e *RateLimitExceededError) Error() string {
	return e.Message
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
		UserRequests: make(map[int]UsersRequests),
		Limit:        limit,
		Interval:     time.Second,
	}
}

func (l *Limiter) Add(id int) error {

	//we can use redis to avoid using map
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
		return &RateLimitExceededError{Message: "Rate limit exceeded"}
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
