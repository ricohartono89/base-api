package way

import "net/http"

// RouterInterface ...
type RouterInterface interface {
	Handle(method string, pattern string, handler http.Handler)
	HandleFunc(method string, pattern string, fn http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
