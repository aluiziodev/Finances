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

func CreateFatura(w http.ResponseWriter, r *http.Request) {

	var req dto.RequestFatura

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("ERRO MULTIPART:", err)
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data_json := r.FormValue("data")
	if err := json.Unmarshal([]byte(data_json), &req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	file, handler, err := r.FormFile("csv")
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	service := getFaturaService()
	id, err := service.CreateFatura(file, handler, req)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Id: criado como sucesso: %s", *id))
}

func ShowFaturas(w http.ResponseWriter, r *http.Request) {
	service := getFaturaService()
	faturas, err := service.GetAllFaturas()
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *faturas)

}

func GetFatura(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	service := getFaturaService()
	fatura, err := service.GetFatura(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *fatura)

}

func DeleteFatura(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getFaturaService()
	err := service.DeleteFatura(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Fatura deletada com sucesso"})
}

func GetFaturaFixo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getFaturaService()
	fatura, err := service.GetFaturaFixo(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *fatura)

}

func GetFaturaParcelado(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getFaturaService()
	fatura, err := service.GetFaturaParcelado(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *fatura)

}

func GetFaturaByCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.RequestCategory

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	service := getBillService()
	fatura, err := service.GetBillsByCategory(id, req.Category)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, *fatura)

}
