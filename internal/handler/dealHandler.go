package handler

import (
	"strconv"
	"net/http"
	"bankDeal/internal/model"
	"github.com/gorilla/schema"
	"github.com/go-playground/validator/v10"
)

type DealHandler struct {
	svc model.DealService
}

func NewDealHandler(svc model.DealService) *DealHandler {
	return &DealHandler{svc: svc}
}


func (h *DealHandler) ListDeals(w http.ResponseWriter, r *http.Request) {
	deals, err := h.svc.ListDeals()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, deals)
}

func (h *DealHandler) GetDeal(w http.ResponseWriter, r *http.Request) {
	
	id, err := strconv.Atoi(r.PathValue("id"))
	
	if err != nil {
		http.Error(w, "無效的 deal ID", http.StatusBadRequest)
		return
	}

	deal, err := h.svc.GetDeal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, deal)
}

func (h *DealHandler) CreateDeal(w http.ResponseWriter, r *http.Request) {

	var requestData model.Deal

    if err := r.ParseForm() ; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	// 將 Form 資料綁定到 Struct
    decoder := schema.NewDecoder()
    if err := decoder.Decode(&requestData, r.PostForm); err != nil {
        http.Error(w, "解析表單失敗", http.StatusBadRequest)
        return
    }


    // 3. 驗證資料合法性
    validate := validator.New()
    if err := validate.Struct(&requestData); err != nil {
        // 這裡可以處理具體的錯誤訊息
        http.Error(w, "驗證失敗: "+err.Error(), http.StatusBadRequest)
        return
    }


	deal, err := h.svc.CreateDeal(requestData.AccountID, requestData.Volume, requestData.TransactionType, requestData.TradingAccountID, requestData.Remark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, deal)
}

