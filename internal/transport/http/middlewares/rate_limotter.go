package middlewares

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Keyer func(*http.Request) string

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor

	rate  rate.Limit
	burst int
	ttl   time.Duration
	key   Keyer

	stop context.CancelFunc
}

func NewRateLimiter(r rate.Limit, burst int, ttl time.Duration, key Keyer) (*RateLimiter, error) {
	if key == nil {
		return nil, errors.New("keyer is nil")
	}
	if r <= 0 {
		return nil, errors.New("rate must be > 0")
	}
	if burst <= 0 {
		return nil, errors.New("burst must be > 0")
	}
	if ttl <= 0 {
		return nil, errors.New("ttl must be > 0")
	}

	ctx, cancel := context.WithCancel(context.Background())

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
		ttl:      ttl,
		key:      key,
		stop:     cancel,
	}

	go rl.cleanupLoop(ctx)
	return rl, nil
}

func (rl *RateLimiter) Close() {
	if rl.stop != nil {
		rl.stop()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k := rl.key(r)
		if k == "" {
			http.Error(w, "bad client key", http.StatusBadRequest)
			return
		}

		lim := rl.getLimiter(k)
		if !lim.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getLimiter(k string) *rate.Limiter {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if v, ok := rl.visitors[k]; ok {
		v.lastSeen = now
		return v.limiter
	}

	lim := rate.NewLimiter(rl.rate, rl.burst)
	rl.visitors[k] = &visitor{limiter: lim, lastSeen: now}
	return lim
}

func (rl *RateLimiter) cleanupLoop(ctx context.Context) {
	interval := rl.ttl / 2
	if interval < time.Second {
		interval = time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := time.Now().Add(-rl.ttl)

			rl.mu.Lock()
			for k, v := range rl.visitors {
				if v.lastSeen.Before(cutoff) {
					delete(rl.visitors, k)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// KeyByRemoteIP — корректно "прямо сейчас" без nginx: берём IP из RemoteAddr (без порта).
func KeyByRemoteIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
