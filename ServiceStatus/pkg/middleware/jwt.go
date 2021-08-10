package middleware

import "net/http"

func Jwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff
		token := r.Header.Get("jwt_token")
		//check jwt_token
		//scriptbot
		r.Header.Set("Name", token) //will change
		next.ServeHTTP(w, r)
	})
}
