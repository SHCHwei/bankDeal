package handler

import (
	"bankDeal/internal/dto/request"
	"bankDeal/internal/model"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type stubUserService struct {
	user          *model.User
	getErr        error
	createErr     error
	createRequest request.CreateUser
}

func (s *stubUserService) GetUser(id int) (*model.User, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.user != nil {
		return s.user, nil
	}
	return nil, errors.New("user not found")
}

func (s *stubUserService) CreateUser(ctx context.Context, requestData request.CreateUser) error {
	s.createRequest = requestData
	return s.createErr
}

func TestUserHandlerCreateUser_Success(t *testing.T) {
	svc := &stubUserService{}
	h := NewUserHandler(svc)

	form := url.Values{}
	form.Set("first_name", "Alice")
	form.Set("last_name", "Lin")
	form.Set("email", "alice@example.com")
	form.Set("phone", "0912345678")
	form.Set("birthdate", "1990-01-01")

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if svc.createRequest.FirstName != "Alice" {
		t.Fatalf("expected first name to be persisted, got %q", svc.createRequest.FirstName)
	}
	if !strings.Contains(rr.Body.String(), "Alice") {
		t.Fatalf("expected response body to contain created user data, got %q", rr.Body.String())
	}
}

func TestUserHandlerGetUser_Success(t *testing.T) {
	svc := &stubUserService{user: &model.User{ID: 1, FirstName: "Alice", LastName: "Lin"}}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	h.GetUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Alice") {
		t.Fatalf("expected body to contain user name, got %q", rr.Body.String())
	}
}

func TestUserHandlerGetUser_InvalidID(t *testing.T) {
	h := NewUserHandler(&stubUserService{})

	req := httptest.NewRequest(http.MethodGet, "/users/not-a-number", nil)
	req.SetPathValue("id", "not-a-number")
	rr := httptest.NewRecorder()

	h.GetUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
