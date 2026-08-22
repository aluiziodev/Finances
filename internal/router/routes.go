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
	{
		URI:      "/fatura/{id}",
		method:   http.MethodDelete,
		function: controllers.DeleteFatura,
	},
	{
		URI:      "/fatura/{id}/parcelado",
		method:   http.MethodGet,
		function: controllers.GetFaturaParcelado,
	},
	{
		URI:      "/fatura/{id}/fixo",
		method:   http.MethodGet,
		function: controllers.GetFaturaFixo,
	},
	{
		URI:      "/fatura/{id}/category",
		method:   http.MethodGet,
		function: controllers.GetFaturaByCategory,
	},
}
