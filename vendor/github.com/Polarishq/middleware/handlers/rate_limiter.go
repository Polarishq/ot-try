package handlers

import (
	"fmt"
	"net/http"

	"github.com/Polarishq/middleware/framework/log"
	"github.com/Polarishq/middleware/framework/ratelimit"
)

// RateLimitHandler runs on each requests to check whether the tenant is within rate
type RateLimitHandler struct {
	handler         http.Handler
	rateLimiter     *ratelimit.RateLimiter
	rateLimitMapper RateLimitMapperFunc
}

// RateLimitMapperFunc exists so that each microservice can define the rate
// for each tier that is passed down by bouncer
// For example, hec-input might allow tier1 tenants to send in 20 events/sec
// but that is a high number for log-query. Log-query might only allow 2 searches/sec
// for a tier1 tenant. This function pushes the responsibility of mapping tiers
// to actual rps (Request Per Second) values to the microservice.
type RateLimitMapperFunc func(r *http.Request) (rps, burst int)

func defaultRateLimitMapper(r *http.Request) (rps, burst int) {
	// The default rate limiter is so high because we want
	// micro-services to override this or fail-safe to (almost) no limiting.
	rps, burst = 100, 20
	return
}

// NewRateLimitHandler creates a new handler that can be added to the middleware stack
func NewRateLimitHandler(handler http.Handler, customRateLimitMapper RateLimitMapperFunc) *RateLimitHandler {
	if customRateLimitMapper == nil {
		customRateLimitMapper = defaultRateLimitMapper
	}
	return &RateLimitHandler{
		handler:         handler,
		rateLimiter:     ratelimit.NewRateLimiter(),
		rateLimitMapper: customRateLimitMapper,
	}
}

// ServeHTTP serves http requests if they are within rate
func (rl *RateLimitHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rps, burst := rl.rateLimitMapper(r)
	limiter := rl.rateLimiter.GetLimiterForTenant(r.Header.Get("X-POLARIS-TENANT-ID"),
		rps, burst)
	if limiter.Allow() {
		rl.handler.ServeHTTP(rw, r)
	} else {
		log.Debugf("Request was rate-limited")
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(429)
		fmt.Fprintln(rw, `{"code": 429, "message": "Rate Exceeded"}`)
	}
}
