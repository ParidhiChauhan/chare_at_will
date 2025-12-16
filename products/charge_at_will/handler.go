package chargeatwill


import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateAuthorization(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Service.CreateAuthorization()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "UPI authorization request initiated",
	})
}
