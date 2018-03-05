package ratelimit

import (
	"github.com/Polarishq/middleware/framework/log"
	"github.com/hashicorp/golang-lru"
	"golang.org/x/time/rate"
)

const cacheSize int = 1024

// RateLimiter is used for limiting the number of requests per second per tenant
type RateLimiter struct {
	cache *lru.ARCCache
}

// NewRateLimiter returns a pointer to a new RateLimiter
func NewRateLimiter() *RateLimiter {
	// TODO: offload this in-memory lru cache to redis
	limitersCache, err := lru.NewARC(cacheSize)
	if err != nil {
		panic("Error creating rate limiter cache: " + err.Error())
	}
	return &RateLimiter{cache: limitersCache}
}

// GetLimiterForTenant gets a token bucket limiter for a given Tenant ID
func (r *RateLimiter) GetLimiterForTenant(tenantID string, rps, burst int) (limiter *rate.Limiter) {
	if limiterIface, ok := r.cache.Get(tenantID); ok {
		limiter, _ = limiterIface.(*rate.Limiter)
	}
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(rps), burst)
		r.cache.Add(tenantID, limiter)
		log.WithFields(log.Fields{
			"rps":      rps,
			"burst":    burst,
			"tenantID": tenantID,
		}).Infof("Creating new rate limiter")
	}
	return
}
