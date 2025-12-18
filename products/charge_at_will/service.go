package chargeatwill

import (
	"bytes"
	"encoding/json"
	"errors"
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
	return &Service{
		key:    key,
		secret: secret,
	}
}

/*
STEP 1:
Create Customer + Order (UPI Mandate)
*/
func (s *Service) CreateCustomerAndOrder(payload string) ([]byte, error) {

	var req map[string]interface{}
	_ = json.Unmarshal([]byte(payload), &req)

	// --------------------
	// 1️⃣ Create Customer
	// --------------------
	customerPayload := map[string]interface{}{
		"name":          req["name"],
		"email":         req["email"],
		"contact":       req["contact"],
		"fail_existing": "0",
	}

	customerBody, _ := json.Marshal(customerPayload)

	customerReq, _ := http.NewRequest(
		"POST",
		"https://api.razorpay.com/v1/customers",
		bytes.NewBuffer(customerBody),
	)

	customerReq.SetBasicAuth(s.key, s.secret)
	customerReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	customerResp, err := client.Do(customerReq)
	if err != nil {
		return nil, err
	}
	defer customerResp.Body.Close()

	customerRespBody, _ := io.ReadAll(customerResp.Body)

	var customerRespJSON map[string]interface{}
	_ = json.Unmarshal(customerRespBody, &customerRespJSON)

	if customerResp.StatusCode >= 400 {
		return nil, errors.New(string(customerRespBody))
	}

	customerID := customerRespJSON["id"].(string)


	orderPayload := map[string]interface{}{
		"amount":      13560, 
		"currency":    "INR",
		"customer_id": customerID,
		"method":      "upi",
		"receipt":     fmt.Sprintf("rcpt_%d", time.Now().Unix()),
		"token": map[string]interface{}{
			"max_amount":      200000,
			"expire_at":       time.Now().AddDate(5, 0, 0).Unix(),
		
			"recurring_value": 1,
		},
	}

	orderBody, _ := json.Marshal(orderPayload)

	orderReq, _ := http.NewRequest(
		"POST",
		"https://api.razorpay.com/v1/orders",
		bytes.NewBuffer(orderBody),
	)

	orderReq.SetBasicAuth(s.key, s.secret)
	orderReq.Header.Set("Content-Type", "application/json")

	orderResp, err := client.Do(orderReq)
	if err != nil {
		return nil, err
	}
	defer orderResp.Body.Close()

	orderRespBody, _ := io.ReadAll(orderResp.Body)

	if orderResp.StatusCode >= 400 {
		return nil, errors.New(string(orderRespBody))
	}

	var orderRespJSON map[string]interface{}
	_ = json.Unmarshal(orderRespBody, &orderRespJSON)

	response := map[string]interface{}{
		"customer_id": customerID,
		"order_id":    orderRespJSON["id"],
	}

	return json.Marshal(response)
}

/*
STEP 2:
Create Authorization Payment (UPI Collect)
*/
func (s *Service) CreateAuthorizationPayment(payload string) ([]byte, error) {

	var req map[string]interface{}
	_ = json.Unmarshal([]byte(payload), &req)

	orderID := req["order_id"].(string)
	amount := int(req["amount"].(float64))

	if orderID == "" {
		return nil, errors.New("order_id missing")
	}

	paymentPayload := map[string]interface{}{
		"amount":   amount, // actual debit amount
		"currency": "INR",
		"order_id": orderID,
		"method":   "upi",
	}

	body, _ := json.Marshal(paymentPayload)

	paymentReq, _ := http.NewRequest(
		"POST",
		"https://api.razorpay.com/v1/payments/create/recurring",
		bytes.NewBuffer(body),
	)

	paymentReq.SetBasicAuth(s.key, s.secret)
	paymentReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(paymentReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(respBody))
	}

	return respBody, nil
}
