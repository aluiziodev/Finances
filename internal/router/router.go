package router

import (
	"net/http"
)

func GenerateRouter() http.Handler {
	mux := http.NewServeMux()
	configure(mux)
	return mux
}

func configure(mux *http.ServeMux) {
	for _, route := range routes {
		mux.HandleFunc(route.method+" "+route.URI, route.function)
	}
}
