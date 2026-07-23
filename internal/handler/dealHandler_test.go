package handler

import (
	"bankDeal/internal/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type stubDealService struct {
	deals            []*model.Deal
	getDealResult    map[string]string
	getDealErr       error
	createDealResult *model.Deal
	createDealErr    error
	accountID        int
	volume           int64
	transactionType  uint8
	tradingAccountID int
	remark           string
}

func (s *stubDealService) ListDeals() ([]*model.Deal, error) {
	return s.deals, nil
}

func (s *stubDealService) GetDeal(id int) (map[string]string, error) {
	if s.getDealErr != nil {
		return nil, s.getDealErr
	}
	if s.getDealResult != nil {
		return s.getDealResult, nil
	}
	return map[string]string{"id": "1"}, nil
}

func (s *stubDealService) CreateDeal(ctx context.Context, accountID int, volume int64, transactionType uint8, tradingAccountID int, remark string) (*model.Deal, error) {
	s.accountID = accountID
	s.volume = volume
	s.transactionType = transactionType
	s.tradingAccountID = tradingAccountID
	s.remark = remark
	return s.createDealResult, s.createDealErr
}

func TestDealHandlerListDeals_Success(t *testing.T) {
	svc := &stubDealService{deals: []*model.Deal{{ID: 1, AccountID: 10}}}
	h := NewDealHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/deals", nil)
	rr := httptest.NewRecorder()

	h.ListDeals(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var deals []*model.Deal
	if err := json.Unmarshal(rr.Body.Bytes(), &deals); err != nil {
		t.Fatalf("expected valid JSON body: %v", err)
	}
	if len(deals) != 1 || deals[0].ID != 1 {
		t.Fatalf("expected one deal in response, got %+v", deals)
	}
}

func TestDealHandlerGetDeal_Success(t *testing.T) {
	svc := &stubDealService{getDealResult: map[string]string{"id": "7"}}
	h := NewDealHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/deals/7", nil)
	req.SetPathValue("id", "7")
	rr := httptest.NewRecorder()

	h.GetDeal(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "7") {
		t.Fatalf("expected response body to contain deal id, got %q", rr.Body.String())
	}
}

func TestDealHandlerGetDeal_InvalidID(t *testing.T) {
	h := NewDealHandler(&stubDealService{})

	req := httptest.NewRequest(http.MethodGet, "/deals/not-a-number", nil)
	req.SetPathValue("id", "not-a-number")
	rr := httptest.NewRecorder()

	h.GetDeal(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDealHandlerCreateDeal_Success(t *testing.T) {
	svc := &stubDealService{createDealResult: &model.Deal{ID: 42}}
	h := NewDealHandler(svc)

	form := url.Values{}
	form.Set("account_id", "1")
	form.Set("volume", "2500")
	form.Set("transaction_type", "1")
	form.Set("trading_account_id", "2")
	form.Set("remark", "test")

	req := httptest.NewRequest(http.MethodPost, "/deals", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.CreateDeal(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if svc.accountID != 1 || svc.volume != 2500 || svc.transactionType != 1 || svc.tradingAccountID != 2 || svc.remark != "test" {
		t.Fatalf("unexpected values passed to service: %+v", svc)
	}
}

func TestDealHandlerCreateDeal_InvalidInput(t *testing.T) {
	h := NewDealHandler(&stubDealService{})

	form := url.Values{}
	form.Set("account_id", "1")
	form.Set("volume", "0")

	req := httptest.NewRequest(http.MethodPost, "/deals", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.CreateDeal(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDealHandlerCreateDeal_ServiceError(t *testing.T) {
	svc := &stubDealService{createDealErr: errors.New("service failed")}
	h := NewDealHandler(svc)

	form := url.Values{}
	form.Set("account_id", "1")
	form.Set("volume", "100")
	form.Set("transaction_type", "1")
	form.Set("trading_account_id", "2")
	form.Set("remark", "fail")

	req := httptest.NewRequest(http.MethodPost, "/deals", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.CreateDeal(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
