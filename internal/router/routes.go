package router

import (
	"finances/internal/controllers"
	"net/http"
)

type route struct {
	URI      string
	method   string
	function func(http.ResponseWriter, *http.Request)
}

var routes = []route{
	{
		URI:      "/fatura",
		method:   http.MethodPost,
		function: controllers.CreateFatura,
	},
	{
		URI:      "/fatura",
		method:   http.MethodGet,
		function: controllers.ShowFaturas,
	},
	{
		URI:      "/fatura/{id}",
		method:   http.MethodGet,
		function: controllers.GetFatura,
	},
}
