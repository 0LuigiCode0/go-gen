package tmp

const MiddlewareTCPTmp = `package {{printf "%v_handler" (index . 0)}}

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func (h *handler) middleware(conn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn.ServeHTTP(w, r)
	})
}

func applyCORS(handler http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	return handlers.CORS(headersOk, originsOk, methodsOk)(handler)
}`
