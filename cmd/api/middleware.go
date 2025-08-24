package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			// built-in recover() to check if a panic occurred.
			pv := recover()
			if pv != nil {

				// Go will close the connection after send the response
				w.Header().Set("Connection", "close")

				// The value returned from recover() has the type any
				// so we have to use fmt.Errorf()
				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	// declare a mutex and a map to hold the client's IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// use the realip.FromRequest() function to get the client's IP address.
		ip := realip.FromRequest(r)

		// lock the mutex to prevent this code from being executed concurrently
		mu.Lock()

		// check to see if the IP address already exists in the map. If it doesn't, then
		// initialize a new rate limiter and add the IP address and limiter on the map
		if _, found := clients[ip]; !found {
			clients[ip] = rate.NewLimiter(2, 4)
		}

		// call the allow() method on the rate limiter for the current IP address.
		// If the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
		// response, just like before.
		if !clients[ip].Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		// very importantly, unlock the mutext before calling the next hander in the chain.
		// Notice that we DON'T use defer to unlock the mutex, as that would mean that the
		// mutex isn't unlocked until all the handlers downstream of this middleware have
		// also returned.
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
