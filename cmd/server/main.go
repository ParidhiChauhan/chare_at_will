package main

import (
	"log"
	"net/http"

	"razorpay-charge-at-will/config"
	"razorpay-charge-at-will/internal/routes"
	chargeatwill "razorpay-charge-at-will/products/charge_at_will"
)

func main() {
	cfg := config.Load()

	service := chargeatwill.NewService(cfg.RzpKey, cfg.RzpSecret)
	handler := chargeatwill.NewHandler(service)

	routes.RegisterChargeAtWillRoutes(handler)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
