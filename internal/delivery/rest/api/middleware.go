package api

import (
	"net/http"
)

func (api *API) maxConcurrent(h http.HandlerFunc, n int) http.HandlerFunc {
	queue := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queue <- struct{}{}
		defer func() { <-queue }()

		h.ServeHTTP(w, r)
	})
}
