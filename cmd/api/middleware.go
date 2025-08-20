package main

import (
	"fmt"
	"net/http"

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

	// Any code here will run only once, when we wrap something with middleware.

	// Initialize a new rate limiter which allows an average of 2 requests per second,
	// with a maximum of 4 requests in a single 'burst'
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Any code here will run for every request that the middleware handles.

		//call limiter.Allow() to see if the request is permitted, and if it's not
		// then we call the rateLimitExceededResponse() helper to return a 429 too many requests
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
