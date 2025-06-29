package ratelimiter

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"tasks/tomlutil"

	"golang.org/x/time/rate"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

// This rate limiter follow the token bucket algorthm to apply rate limiting

/*
RateLimitInitialBurstReqCount ==> burst size (bucket capacity).
RateLimitReqPerSecond ==> refill rate (tokens per second).

*/

func AssignRateLimitValue() *rate.Limiter {
	// Reading TOML values.
	lTeXlConfig := tomlutil.ReadTomlConfig("toml/config.toml")

	lTRateLimitReqPerSecond := fmt.Sprintf("%v", lTeXlConfig.(map[string]any)["RateLimitReqPerSecond"])
	lTRateLimitInitialBurstReqCount := fmt.Sprintf("%v", lTeXlConfig.(map[string]any)["RateLimitInitialBurstReqCount"])

	lProceedValue, _ := strconv.ParseFloat(lTRateLimitReqPerSecond, 64)
	rAcceptValue, _ := strconv.Atoi(lTRateLimitInitialBurstReqCount)

	rProceedValue := rate.Limit(lProceedValue)
	lLimiter := rate.NewLimiter(rProceedValue, rAcceptValue)

	return lLimiter
}

/* func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {

	TokenBucketCapacity, TokenRefreshRate := AssignRateLimitValue()
	lLimiter := rate.NewLimiter(TokenBucketCapacity, TokenRefreshRate)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lLimiter.Allow() {
			lMessage := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&lMessage)
			return
		} else {
			next(w, r)
		}
	})
} */

// func RateLimiter(next http.Handler) http.Handler {
func RateLimiterPerclient(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	// Mutex to handle concurrent access to the map
	var lMutex sync.Mutex
	// Map to store rate limiters per client
	lClients := make(map[string]*rate.Limiter)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// clientIP := r.RemoteAddr
		// You’re using r.RemoteAddr to identify the client. This includes both IP and port.

		lClientIP, _, _ := net.SplitHostPort(r.RemoteAddr) // Strip the port

		lMutex.Lock()
		lLimiter, lExists := lClients[lClientIP]
		if !lExists {
			// Create a new rate limiter for a client
			lLimiter = rate.NewLimiter(5, 5) // Allows 2 requests immediately, then 1 request per second
			lClients[lClientIP] = lLimiter
		}
		lMutex.Unlock()

		// Limit the number of requests
		if !lLimiter.Allow() {
			lMessage := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, try again later.",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&lMessage)
			return
		} else {
			next(w, r)
		}
	})
}
