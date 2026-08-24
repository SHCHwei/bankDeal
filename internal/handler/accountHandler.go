package handler

import (
	"bankDeal/internal/dto/request"
	"bankDeal/internal/model"
	"net/http"
	"strconv"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type AccountHandler struct {
	svc 		model.AccountService
	decode 		*schema.Decoder
	validate 	*validator.Validate
}

func NewAccountHandler(svc model.AccountService) *AccountHandler {
	return &AccountHandler{
		svc: svc,
		decode: schema.NewDecoder(),
		validate: validator.New(),
	}
}


func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	
	id, err := strconv.Atoi(r.PathValue("id"))

	
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	
	account, err := h.svc.GetAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, account)
}