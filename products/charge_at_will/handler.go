package chargeatwill

import (
	"encoding/json"
	"net/http"
)

/* ---------------- STEP 1 & 2: CUSTOMER + ORDER ---------------- */

func (h *Handler) CreateAuthorization(w http.ResponseWriter, r *http.Request) {
	var req SetupRequest
	json.NewDecoder(r.Body).Decode(&req)

	customerID, err := h.service.CreateCustomer(&CreateCustomerRequest{
		Name:    req.Name,
		Email:   req.Email,
		Contact: req.Contact,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	orderID, err := h.service.CreateOrder(customerID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"customer_id": customerID,
		"order_id":    orderID,
	})
}

/* ---------------- STEP 3: AUTHORIZATION PAYMENT (MISSING STEP) ---------------- */

func (h *Handler) CreateAuthorizationPayment(w http.ResponseWriter, r *http.Request) {
	var req AuthorizationPaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.service.CreateAuthorizationPayment(req.OrderID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
