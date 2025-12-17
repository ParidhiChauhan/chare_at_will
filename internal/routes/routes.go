package routes

import (
	"net/http"

	chargeatwill "razorpay_charge_at_will/products/charge_at_will"
)

func RegisterChargeAtWillRoutes(h *chargeatwill.Handler) {
	http.HandleFunc("/charge-at-will/authorize", h.CreateAuthorization)
}
