package controllers

import (
	"encoding/json"
	"finances/internal/dto"
	"finances/internal/response"
	"finances/internal/service"
	"fmt"
	"net/http"
)

var getFaturaService = func() *service.FaturaService {
	return service.NewFaturaService()
}
var getBillService = func() *service.BillService {
	return service.NewBillService()
}

var getCardService = func() *service.CardService {
	return service.NewCardService()
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestCard
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	service := getCardService()
	id, err := service.CreateCard(req)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Id: criado como sucesso: %s", *id))
}

func ShowCards(w http.ResponseWriter, r *http.Request) {
	service := getCardService()
	cards, err := service.GetAllCards()
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *cards)
}

func GetCardById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getCardService()
	card, err := service.GetCardById(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *card)
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getCardService()
	err := service.DeleteCard(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, fmt.Sprintf("Cartão deletado com sucesso: %s", id))
}
