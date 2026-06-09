package handler


import (
	"net/http"
	"strconv"
	"bankDeal/internal/model"
	"bankDeal/internal/dto/request"
	"github.com/gorilla/schema"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	svc model.UserService
}


func NewUserHandler(svc model.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}


func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	
	var requestData request.CreateUser


    err := r.ParseForm()
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	// 將 Form 資料綁定到 Struct
    decoder := schema.NewDecoder()
    if err := decoder.Decode(&requestData, r.PostForm); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// 驗證資料合法性
    validate := validator.New()
    if err := validate.Struct(&requestData); err != nil {
        // 這裡可以處理具體的錯誤訊息
        http.Error(w, "驗證失敗: "+err.Error(), http.StatusBadRequest)
        return
    }
	

	err = h.svc.CreateUser(requestData)

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
		http.Error(w, "user id is not found", http.StatusBadRequest)
		return
	}


	writeJSON(w, http.StatusCreated, user)
}



