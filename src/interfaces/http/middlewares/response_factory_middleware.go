package middlewares

import (
	"net/http"

	httpserver "base/src/infrastructure/http_server"
)

func ResponseFactory(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpResponseFactory := httpserver.NewHttpResponseFactory(w)
		handler.ServeHTTP(httpResponseFactory, r)
	}
}
