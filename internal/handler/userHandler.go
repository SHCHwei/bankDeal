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

type UserHandler struct {
	svc 		model.UserService
	decode 		*schema.Decoder
	validate 	*validator.Validate
}

func NewUserHandler(svc model.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
		decode: schema.NewDecoder(),
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var requestData request.CreateUser

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 將 Form 資料綁定到 Struct
	if err := h.decode.Decode(&requestData, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 驗證資料合法性
	if err := h.validate.Struct(&requestData); err != nil {
		// 這裡可以處理具體的錯誤訊息
		http.Error(w, "驗證失敗: " + err.Error(), http.StatusBadRequest)
		return
	}

	err = h.svc.CreateUser(context.Background(), requestData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, requestData)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))


	if err != nil {
		http.Error(w, "無效的 user ID", http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
