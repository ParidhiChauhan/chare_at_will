package main

import (
	"log"
	"net/http"

	"razorpay_charge_at_will/config"
	"razorpay_charge_at_will/internal/routes"
	chargeatwill "razorpay_charge_at_will/products/charge_at_will"
)

func main() {
	cfg := config.Load()

	service := chargeatwill.NewService(cfg.RzpKey, cfg.RzpSecret)
	handler := chargeatwill.NewHandler(service)

	routes.RegisterChargeAtWillRoutes(handler)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
