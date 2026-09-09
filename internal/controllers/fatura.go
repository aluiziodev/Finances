package controllers

import (
	"encoding/json"
	"finances/internal/dto"
	"finances/internal/response"
	"fmt"
	"net/http"
)

func CreateFatura(w http.ResponseWriter, r *http.Request) {
	card_id := r.PathValue("card_id")
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
	id, err := service.CreateFatura(file, handler, req, card_id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, fmt.Sprintf("Id: criado como sucesso: %s", *id))
}

func ShowFaturas(w http.ResponseWriter, r *http.Request) {
	card_id := r.PathValue("card_id")
	service := getFaturaService()
	faturas, err := service.GetAllFaturas(card_id)
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

	service := getBillService()
	fatura, err := service.GetFaturaFixo(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *fatura)

}

func GetFaturaParcelado(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getBillService()
	fatura, err := service.GetFaturaParcelado(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, *fatura)

}

func GetFaturaByCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	category := r.URL.Query().Get("category")
	if category == "" {
		response.ErrorResponse(w, http.StatusBadRequest, "missing category query parameter")
		return
	}

	service := getBillService()
	fatura, err := service.GetBillsByCategory(id, category)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, *fatura)

}

func UpdateFaturaPaid(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getFaturaService()
	err := service.UpdateFaturaPaid(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Status atualizado com sucesso"})
}

func UpdateFaturaPending(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := getFaturaService()
	err := service.UpdateFaturaPending(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "Status atualizado com sucesso"})
}
