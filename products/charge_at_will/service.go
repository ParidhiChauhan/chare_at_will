package chargeatwill

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	key    string
	secret string
}

func NewService(key, secret string) *Service {
	return &Service{key: key, secret: secret}
}

/* ---------------- CREATE CUSTOMER ---------------- */

func (s *Service) CreateCustomer(req *CreateCustomerRequest) (string, error) {
	url := "https://api.razorpay.com/v1/customers"

	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.SetBasicAuth(s.key, s.secret)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("customer error: %s", string(respBody))
	}

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	return result["id"].(string), nil
}

/* ---------------- CREATE ORDER ---------------- */

func (s *Service) CreateOrder(customerID string, amount int) (string, error) {
	url := "https://api.razorpay.com/v1/orders"

	payload := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"customer_id": customerID,
		"method":   "upi",
		"receipt":  fmt.Sprintf("receipt_%d", time.Now().Unix()),
		"token": map[string]interface{}{
			"max_amount":     200000,
			"expire_at":      time.Now().AddDate(1, 0, 0).Unix(),
			"frequency":      "as_presented",
			"recurring_type": "as_presented",
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.SetBasicAuth(s.key, s.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("order error: %s", string(respBody))
	}

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	return result["id"].(string), nil
}

/* ---------------- CREATE AUTHORIZATION PAYMENT (NEW & REQUIRED) ---------------- */

func (s *Service) CreateAuthorizationPayment(orderID string) ([]byte, error) {
	url := "https://api.razorpay.com/v1/payments/create/authorization"

	payload := map[string]string{
		"order_id": orderID,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.SetBasicAuth(s.key, s.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("authorization error: %s", string(respBody))
	}

	return respBody, nil
}
