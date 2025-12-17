package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)

// --- 1. Structs for JSON Payloads ---

// OrderRequest represents the body for Creating an Order
type OrderRequest struct {
    Amount         int                    `json:"amount"`
    Currency       string                 `json:"currency"`
    PaymentCapture bool                   `json:"payment_capture"`
    Receipt        string                 `json:"receipt"`
    Notification   map[string]interface{} `json:"notification"`
    Notes          map[string]string      `json:"notes"`
}

// OrderResponse represents the response from the Order API
type OrderResponse struct {
    ID        string            `json:"id"`
    Entity    string            `json:"entity"`
    Amount    int               `json:"amount"`
    Status    string            `json:"status"`
    Notes     map[string]string `json:"notes"`
    CreatedAt int64             `json:"created_at"`
    // Add other fields if needed, strictly we need ID for the next step
}

// RecurringPaymentRequest represents the body for Creating a Recurring Payment
type RecurringPaymentRequest struct {
    Email       string            `json:"email"`
    Contact     string            `json:"contact"`
    Amount      int               `json:"amount"`
    Currency    string            `json:"currency"`
    OrderID     string            `json:"order_id"`
    CustomerID  string            `json:"customer_id"`
    Token       string            `json:"token"`
    Recurring   bool              `json:"recurring"`
    Description string            `json:"description"`
    Notes       map[string]string `json:"notes"`
}

// PaymentResponse represents the response from the Payment API
type PaymentResponse struct {
    RazorpayPaymentID string `json:"razorpay_payment_id"`
}

// --- 2. Database Simulation ---

// GetCredentialsMock mimics fetching static keys from a DB
func GetCredentialsMock() (string, string) {
    // TODO: Replace these with actual Razorpay Key ID and Secret
    keyID := "rzp_test_Df9PS5NyT01KEo"
    keySecret := "OcyZcfdwniBefW7tnprMc4mI"
    return keyID, keySecret
}

// --- 3. HTTP Client Helper ---

// sendRazorpayRequest handles the HTTP boilerplate, Basic Auth, and JSON Marshaling
func sendRazorpayRequest(method, url string, payload interface{}, target interface{}, keyId, keySecret string) error {
    client := &http.Client{Timeout: 10 * time.Second}

    // Marshal payload to JSON
    jsonBytes, err := json.Marshal(payload)
    if err != nil {
       return fmt.Errorf("error marshaling payload: %v", err)
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
    if err != nil {
       return fmt.Errorf("error creating request: %v", err)
    }

    // Set Headers and Auth
    req.SetBasicAuth(keyId, keySecret)
    req.Header.Set("Content-Type", "application/json")

    // Execute Request
    resp, err := client.Do(req)
    if err != nil {
       return fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    // Read Body
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
       return fmt.Errorf("error reading response body: %v", err)
    }

    // Print Raw Response for the terminal output
    fmt.Printf("\n--- Response from %s ---\n", url)
    prettyPrintJSON(bodyBytes)

    // Check for non-200 status codes
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
       return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
    }

    // Unmarshal into target struct to use data programmatically
    return json.Unmarshal(bodyBytes, target)
}

// prettyPrintJSON is a utility to print the JSON nicely to the terminal
func prettyPrintJSON(data []byte) {
    var prettyJSON bytes.Buffer
    if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
       fmt.Println(string(data))
    } else {
       fmt.Println(prettyJSON.String())
    }
}

// --- 4. Main Execution Flow ---

func main() {
    // 1. Get Credentials
    keyID, keySecret := GetCredentialsMock()
    fmt.Println("Credentials fetched from DB...")

    // ------------------------------------------
    // STEP 1: Create Order
    // ------------------------------------------
    fmt.Println("\nExecuting Step 1: Creating Order...")

    orderPayload := OrderRequest{
       Amount:         1000, // Amount in paise (1000 = 10.00)
       Currency:       "INR",
       PaymentCapture: true,
       Receipt:        "Receipt No. 190",
       Notification: map[string]interface{}{
          "token_id":      "token_RsEEEhIWuNaEVw", // Ensure this token exists in your Sandbox
          "payment_after": 1634057114,
       },
       Notes: map[string]string{
          "notes_key_1": "Tea, Earl Grey, Hot",
          "notes_key_2": "Tea, Earl Grey… decaf.",
       },
    }

    var orderResp OrderResponse
    err := sendRazorpayRequest("POST", "https://api.razorpay.com/v1/orders", orderPayload, &orderResp, keyID, keySecret)
    if err != nil {
       log.Fatalf("Step 1 Failed: %v", err)
    }

    // ------------------------------------------
    // STEP 2: Create Recurring Payment
    // ------------------------------------------
    // We use the ID received from Step 1
    fmt.Printf("\nExecuting Step 2: Creating Recurring Payment using Order ID: %s...\n", orderResp.ID)

    paymentPayload := RecurringPaymentRequest{
       Email:       "gaurav.kumar@example.com",
       Contact:     "9000090000",
       Amount:      1000,
       Currency:    "INR",
       OrderID:     orderResp.ID,           // DYNAMIC: From Step 1
       CustomerID:  "cust_R9SwX5UQqXcjOn",  // TODO: Replace with valid Customer ID
       Token:       "token_RsEEEhIWuNaEVw", // TODO: Replace with valid Token ID
       Recurring:   true,
       Description: "Creating recurring payment for Gaurav Kumar",
       Notes: map[string]string{
          "note_key 1": "Beam me up Scotty",
          "note_key 2": "Tea. Earl Gray. Hot.",
       },
    }

    var paymentResp PaymentResponse
    // Note: The URL provided in your prompt implies a direct payment creation.
    // Ensure this endpoint is valid for your specific integration flow.
    err = sendRazorpayRequest("POST", "https://api.razorpay.com/v1/payments/create/recurring/", paymentPayload, &paymentResp, keyID, keySecret)
    if err != nil {
       log.Fatalf("Step 2 Failed: %v", err)
    }

    fmt.Println("\n--- Workflow Completed Successfully ---")
}










