package handler

import (
	"bankDeal/internal/model"
	"bankDeal/internal/dto/request"
	"bankDeal/internal/logger"
	"net/http"
	"strconv"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type DealHandler struct {
	svc 		model.DealService
   	decoder		*schema.Decoder
    validate	*validator.Validate
}

func NewDealHandler(svc model.DealService) *DealHandler {
	return &DealHandler{
		svc: svc,
		decoder:  schema.NewDecoder(),
        validate: validator.New(),
	}
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

	var requestData request.CreateDeal

	// 日誌紀錄
	log, err := logger.GetInstance()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.LogRequestError(0, 0, 0, 0, "ParseForm", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 將 Form 資料綁定到 Struct
	if err := h.decoder.Decode(&requestData, r.PostForm); err != nil {
		log.LogRequestError(0, 0, 0, 0, "Decode", "解析表單失敗")
		http.Error(w, "解析表單失敗", http.StatusBadRequest)
		return
	}

	// 3. 驗證資料合法性
	if err := h.validate.Struct(&requestData); err != nil {
		// 記錄驗證錯誤
		log.LogRequestError(
			requestData.AccountID,
			requestData.Volume,
			requestData.TransactionType,
			requestData.TradingAccountID,
			"Validation",
			err.Error(),
		)
		http.Error(w, "驗證失敗: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 4. 執行交易
	deal, err := h.svc.CreateDeal(
		context.Background(),
		requestData.AccountID,
		requestData.Volume,
		requestData.TransactionType,
		requestData.TradingAccountID,
		requestData.Remark,
	)
	
	if err != nil {
		// 記錄交易創建失敗
		log.LogCreateDeal(
			requestData.AccountID,
			requestData.Volume,
			requestData.TransactionType,
			requestData.TradingAccountID,
			requestData.Remark,
			false,
			0,
			err.Error(),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 記錄成功的交易創建
	log.LogCreateDeal(
		requestData.AccountID,
		requestData.Volume,
		requestData.TransactionType,
		requestData.TradingAccountID,
		requestData.Remark,
		true,
		deal.ID,
		"",
	)

	writeJSON(w, http.StatusCreated, deal)
}
