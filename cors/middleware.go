package cors

import (
	"net/http"
)

func MiddlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type")
		handler.ServeHTTP(writer, request)
	})
}
