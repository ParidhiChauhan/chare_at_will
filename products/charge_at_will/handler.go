package chargeatwill

import (
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

/*
STEP 1:
Create customer + order
*/
func (h *Handler) CreateAuthorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)

	resp, err := h.service.CreateCustomerAndOrder(string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("order error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

/*
STEP 2:
Create authorization payment (UPI Autopay approval)
*/
func (h *Handler) CreateAuthorizationPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)

	resp, err := h.service.CreateAuthorizationPayment(string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("authorization error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
