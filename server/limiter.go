package main

import (
	"log"
	"net/http"
	"net/netip"
	"os"
	"sync"
	"time"
)

const rateLimitWindow = time.Minute * 15
const rateLimitRequests = 1000

type rateLimitClient struct {
	WindowStart  time.Time
	RequestCount int
}

// basic rate limiter for the time being
type rateLimiter struct {
	clients map[string]rateLimitClient
	mut     sync.RWMutex
}

// clean removes old entries from clients tracked by the rate limiter.
func (rl *rateLimiter) clean() {
	rl.mut.Lock()
	defer rl.mut.Unlock()

	for ip, e := range rl.clients {
		if e.WindowStart.Add(rateLimitWindow).Before(time.Now()) {
			delete(rl.clients, ip)
		}
	}
}

// Limit returns true if a request should be rate-limited.
//
// If true, it responds to the request with { "error": "Too many requests" }
// and status code 429.
func (rl *rateLimiter) Limit(res http.ResponseWriter, req *http.Request) bool {
	var ip string
	ipHeader := os.Getenv("IP_HEADER")
	if ipHeader == "" {
		addrport, err := netip.ParseAddrPort(req.RemoteAddr)
		if err != nil {
			writeErr(res, "Too many requests", 429)
			log.Printf("error parsing remote address: %v", err)
			return true
		}
		ip = addrport.Addr().String()
	} else {
		addr, err := netip.ParseAddr(ipHeader)
		log.Printf("error parsing address: %v", err)
		ip = addr.String()
	}

	rl.mut.RLock()
	client, ok := rl.clients[ip]
	rl.mut.RUnlock()
	if !ok {
		rl.clean()
		client.WindowStart = time.Now()
		client.RequestCount = 1
	} else {
		if client.WindowStart.Add(rateLimitWindow).Before(time.Now()) {
			client.WindowStart = time.Now()
			client.RequestCount = 1
		} else {
			client.RequestCount += 1
		}
	}

	if client.RequestCount > rateLimitRequests {
		writeErr(res, "Too many requests", 429)
		return true
	}

	rl.mut.Lock()
	defer rl.mut.Unlock()
	rl.clients[ip] = client
	return false
}
