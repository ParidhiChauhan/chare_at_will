package main

import (
	"log"
	"net/http"
	

	"razorpay_charge_at_will/config"
	"razorpay_charge_at_will/internal/routes"
	chargeatwill "razorpay_charge_at_will/products/charge_at_will"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize Razorpay service & handler
	service := chargeatwill.NewService(cfg.RzpKey, cfg.RzpSecret)
	handler := chargeatwill.NewHandler(service)

	// Register backend API routes
	routes.RegisterChargeAtWillRoutes(handler)

	// Serve frontend files from "frontend/" folder
	frontendDir := "./frontend"
	fs := http.FileServer(http.Dir(frontendDir))
	// Handle "/" to serve index.html by default
	http.Handle("/", fs)

	// Optional: if you want API and frontend under same mux
	// http.Handle("/charge-at-will/", http.StripPrefix("/charge-at-will/", handler))

	// Start server
	log.Println("Server running on port", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
