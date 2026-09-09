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
		URI:      "/card",
		method:   http.MethodPost,
		function: controllers.CreateCard,
	},
	{
		URI:      "/card",
		method:   http.MethodGet,
		function: controllers.ShowCards,
	},
	{
		URI:      "/card/{id}",
		method:   http.MethodGet,
		function: controllers.GetCardById,
	},
	{
		URI:      "/card/{id}",
		method:   http.MethodDelete,
		function: controllers.DeleteCard,
	},
	{
		URI:      "/card/{card_id}/fatura",
		method:   http.MethodPost,
		function: controllers.CreateFatura,
	},
	{
		URI:      "/card/{card_id}/fatura",
		method:   http.MethodGet,
		function: controllers.ShowFaturas,
	},
	{
		URI:      "/card/{card_id}/fatura/{id}",
		method:   http.MethodGet,
		function: controllers.GetFatura,
	},
	{
		URI:      "/card/{card_id}/fatura/{id}",
		method:   http.MethodDelete,
		function: controllers.DeleteFatura,
	},
	{
		URI:      "/card/{card_id}/fatura/{id}/parcelado",
		method:   http.MethodGet,
		function: controllers.GetFaturaParcelado,
	},
	{
		URI:      "/card/{card_id}/fatura/{id}/fixo",
		method:   http.MethodGet,
		function: controllers.GetFaturaFixo,
	},
	{
		URI:      "/card/{card_id}/fatura/{id}/category/",
		method:   http.MethodGet,
		function: controllers.GetFaturaByCategory,
	},
}
