package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/grglucastr/go-greenlight/internal/data"
	"github.com/grglucastr/go-greenlight/internal/validator"
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

	// if rate limiting is not enabled
	if !app.config.limiter.enabled {
		return next
	}

	// Client struct to hold the rate limiter and last seen time for each
	// client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// declare a mutex and a map to hold the client's IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background goroutine which removes old entries from the clients map once
	// every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			// lock the mutex to prevent any reate limiter checks from happening while
			// the cleanup is taking place
			mu.Lock()

			// loop through all clients. If they haven't been seen within the last three
			// minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			// importantly, unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// use the realip.FromRequest() function to get the client's IP address.
		ip := realip.FromRequest(r)

		// lock the mutex to prevent this code from being executed concurrently
		mu.Lock()

		// check to see if the IP address already exists in the map. If it doesn't, then
		// initialize a new rate limiter and add the IP address and limiter on the map
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
			}
		}

		clients[ip].lastSeen = time.Now()

		// call the allow() method on the rate limiter for the current IP address.
		// If the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
		// response, just like before.
		if !clients[ip].limiter.Allow() {
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Indicates to any caches that the response may vary
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
