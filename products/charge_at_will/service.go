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

func (s *Service) doRequest(method, url string, body interface{}) ([]byte, error) {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.SetBasicAuth(s.key, s.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("razorpay API error: %s", respBody)
	}

	return respBody, nil
}

/* STEP 1: CREATE CUSTOMER */
func (s *Service) CreateCustomer(req *AuthorizationRequest) (string, error) {
	payload := map[string]interface{}{
		"name":          req.Name,
		"email":         req.Email,
		"contact":       req.Contact,
		"fail_existing": "0",
		"notes": map[string]string{
			"note_key_1": "September",
			"note_key_2": "Make it so.",
		},
	}

	resp, err := s.doRequest("POST", RazorpayBaseURL+"/customers", payload)
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	json.Unmarshal(resp, &result)

	return result.ID, nil
}

/* STEP 2: CREATE ORDER (UPI + TOKEN) */
func (s *Service) CreateOrder(req *AuthorizationRequest, customerID string) (string, error) {
	payload := map[string]interface{}{
		"amount":      req.Amount,
		"currency":    req.Currency,
		"customer_id": customerID,
		"method":      "upi",

		// ✅ UNIQUE RECEIPT (FIX)
		"receipt": fmt.Sprintf("receipt_%d", time.Now().UnixNano()),

		"token": map[string]interface{}{
			"max_amount":      200000,
			"expire_at":       2709971120,
			"frequency":       "as_presented",
		
		},

		"notes": map[string]string{
			"notes_key_1": "Tea, Earl Grey, Hot",
			"notes_key_2": "Tea, Earl Grey… decaf.",
		},
	}

	resp, err := s.doRequest("POST", RazorpayBaseURL+"/orders", payload)
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	json.Unmarshal(resp, &result)

	return result.ID, nil
}
