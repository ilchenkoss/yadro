package utils

import (
	"myapp/internal-xkcd/config"
	"myapp/internal-xkcd/core/domain"
	"sync"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	httpCfg := &config.HttpServerConfig{
		ConcurrencyLimit: 2,
		RateLimit:        10,
	}

	limiter := NewLimiter(httpCfg)

	if limiter.cl == nil || limiter.rl == nil {
		t.Errorf("Expected Limiter to initialize ConcurrencyLimiter and RateLimiter")
	}
}

func TestConcurrencyLimiter(t *testing.T) {
	cl := NewConcurrencyLimiter(1)

	var wg sync.WaitGroup
	wg.Add(2)

	start := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer wg.Done()
		<-start
		cl.Add()
		time.Sleep(100 * time.Millisecond) //simulate work
		cl.Done()
		done <- struct{}{}
	}()

	go func() {
		defer wg.Done()
		<-start
		cl.Add()
		cl.Done()
		done <- struct{}{}
	}()

	close(start)
	<-done
	<-done

	wg.Wait()
}

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5)

	userID := uint64(1)

	//simulate 5 requests
	for i := 0; i < 5; i++ {
		err := rl.Add(userID)
		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i+1, err)
		}
	}

	//6th request should fail
	err := rl.Add(userID)
	if err != domain.ErrRateLimitExceeded {
		t.Errorf("Expected ErrRateLimitExceeded, got %v", err)
	}

	//wait interval and test again
	time.Sleep(rl.Interval + 10*time.Millisecond)

	err = rl.Add(userID)
	if err != nil {
		t.Errorf("Unexpected error after interval: %v", err)
	}
}
