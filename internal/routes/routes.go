package routes

import (
	"net/http"

	chargeatwill "razorpay-charge-at-will/products/charge_at_will"
)

func RegisterChargeAtWillRoutes(handler *chargeatwill.Handler) {
	http.HandleFunc("/charge-at-will/authorize", handler.CreateAuthorization)
}
