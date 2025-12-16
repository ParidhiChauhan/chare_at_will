package chargeatwill

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Service struct {
	KeyID     string
	KeySecret string
}

func NewService(key, secret string) *Service {
	return &Service{
		KeyID:     key,
		KeySecret: secret,
	}
}

func (s *Service) CreateAuthorization() (*http.Response, error) {
	payload := map[string]interface{}{
		"amount":   100,
		"currency": "INR",
		"method":   "upi",
		"customer": map[string]string{
			"email":   "customer@example.com",
			"contact": "9999999999",
		},
		"upi": map[string]interface{}{
			"flow": "collect",
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", RazorpayAuthURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.KeyID, s.KeySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
