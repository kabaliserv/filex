package files

import "net/http"

func getFileMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func FileMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
