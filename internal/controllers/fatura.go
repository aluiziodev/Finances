package controllers

import (
	"encoding/json"
	"finances/internal/models"
	"finances/internal/parser"
	"finances/internal/response"
	"net/http"
)

func GetFatura(w http.ResponseWriter, r *http.Request) {

	var req models.RequestFatura
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Erro no corpo da requisiçao!")
		return
	}

	fatura, err := parser.ParserCSVtoModels(req.Path, req.Description, req.Status)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Erro no parser do csv!")
		return
	}

	response.WriteJSON(w, http.StatusOK, *fatura)
}
