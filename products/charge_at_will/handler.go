package chargeatwill

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAuthorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthorizationRequest
	json.NewDecoder(r.Body).Decode(&req)

	customerID, err := h.service.CreateCustomer(&req)
	if err != nil {
		http.Error(w, "Customer error: "+err.Error(), 500)
		return
	}

	orderID, err := h.service.CreateOrder(&req, customerID)
	if err != nil {
		http.Error(w, "Order error: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"customer_id": customerID,
		"order_id":    orderID,
	})
}
