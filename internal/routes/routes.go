package routes

import (
	"net/http"

	chargeatwill "razorpay_charge_at_will/products/charge_at_will"
)

func RegisterChargeAtWillRoutes(handler *chargeatwill.Handler) {
	http.HandleFunc("/charge-at-will/authorize", handler.CreateAuthorization)
	http.HandleFunc("/charge-at-will/authorize-payment", handler.CreateAuthorizationPayment)
}
