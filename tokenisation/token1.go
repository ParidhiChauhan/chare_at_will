package main
import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)
// --- 1. Structs for JSON Payloads & Responses ---
// PaymentResponse represents the response from GET /payments/:id
type PaymentResponse struct {
    ID         string            `json:"id"`
    Entity     string            `json:"entity"`
    Amount     int               `json:"amount"`
    Status     string            `json:"status"`
    CustomerID string            `json:"customer_id"`
    TokenID    string            `json:"token_id"` // This is the ID used for recurring payments
    Email      string            `json:"email"`
    Contact    string            `json:"contact"`
    CreatedAt  int64             `json:"created_at"`
    Notes      map[string]string `json:"notes"`
}
// --- 2. Configuration ---
const (
    RAZORPAY_KEY_ID     = "rzp_test_Df9PS5NyT01KEo"
    RAZORPAY_KEY_SECRET = "OcyZcfdwniBefW7tnprMc4mI"
    RAZORPAY_API_BASE   = "https://api.razorpay.com/v1"
)
// --- 3. HTTP Client Helper ---
func sendRazorpayRequest(method, url string, target interface{}) error {
    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest(method, url, nil)
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }
    // Set Headers and Auth
    req.SetBasicAuth(RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET)
    req.Header.Set("Content-Type", "application/json")
    // Execute Request
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("error reading response: %v", err)
    }
    // Check for non-200 status codes
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("API failed with status %d: %s", resp.StatusCode, string(bodyBytes))
    }
    // Unmarshal into target struct
    return json.Unmarshal(bodyBytes, target)
}
// --- 4. Main Execution ---
func main() {
    // The ID of the initial payment where card saving consent was given
    paymentID := "pay_RsEEDMcc2SPXVk"
    fmt.Printf("Step 1: Fetching Token details using Payment ID: %s...\n", paymentID)
    // URL: GET https://api.razorpay.com/v1/payments/{payment_id}
    fetchURL := fmt.Sprintf("%s/payments/%s", RAZORPAY_API_BASE, paymentID)
    var paymentInfo PaymentResponse
    err := sendRazorpayRequest("GET", fetchURL, &paymentInfo)
    if err != nil {
        log.Fatalf(":x: Step 1 Failed: %v", err)
    }
    // Output the results
    fmt.Println("\n--- Payment Details Fetched ---")
    fmt.Printf("Payment ID:   %s\n", paymentInfo.ID)
    fmt.Printf("Status:       %s\n", paymentInfo.Status)
    fmt.Printf("Customer ID:  %s\n", paymentInfo.CustomerID)
    if paymentInfo.TokenID != "" {
        fmt.Println("\n==============================================")
        fmt.Printf(":fire: SUCCESSFULLY EXTRACTED TOKEN ID: %s\n", paymentInfo.TokenID)
        fmt.Println("==============================================")
        fmt.Println("You can now use this token to execute 'Charge at Will' payments.")
    } else {
        fmt.Println("\n:warning:  Token ID was not found for this payment.")
        fmt.Println("Ensure the original payment was created with 'recurring: true'.")
    }
}
