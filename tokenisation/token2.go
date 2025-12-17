package main
import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)
// !!! IMPORTANT: REPLACE THESE CONSTANTS WITH YOUR ACTUAL VALUES !!!
const (
    // Your Razorpay Key ID (rzp_test_...)
    RAZORPAY_KEY_ID    = "rzp_test_Df9PS5NyT01KEo"
    // Your Razorpay Key Secret (The actual key)
    RAZORPAY_KEY_SECRET = "OcyZcfdwniBefW7tnprMc4mI"
    // The unique identifier of the customer (cust_...).
    // This customer must have previously saved a card (token).
    CUSTOMER_ID        = "cust_GcwirOxeZwg63m"
    // API endpoint format string: /v1/customers/{CUSTOMER_ID}/tokens
    RAZORPAY_API_URL   = "https://api.razorpay.com/v1/customers/%s/tokens"
)
func main() {
    // 1. Construct the full API URL
    // Example: https://api.razorpay.com/v1/customers/cust_GcwirOxeZwg63m/tokens
    fetchURL := fmt.Sprintf(RAZORPAY_API_URL, CUSTOMER_ID)
    // 2. Create the HTTP request (GET method)
    req, err := http.NewRequest("GET", fetchURL, nil)
    if err != nil {
        log.Fatalf("Error creating HTTP request: %v", err)
    }
    // 3. Set Basic Authentication: Key ID is the username, Key Secret is the password
    req.SetBasicAuth(RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET)
    // 4. Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error sending request to Razorpay: %v", err)
    }
    // Ensure the response body is closed when the function returns
    defer resp.Body.Close()
    // 5. Read and parse the response body
    body, _ := io.ReadAll(resp.Body)
    // Handle non-200 status codes (like 400 Bad Request)
    if resp.StatusCode != http.StatusOK {
        log.Fatalf(":x: API Request failed. Status: %s. Response: %s", resp.Status, body)
    }
    // Razorpay returns a collection entity
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        log.Fatalf("Error unmarshalling response JSON: %v", err)
    }
    // 6. Extract the list of tokens (which is inside the 'items' key)
    items, ok := result["items"].([]interface{})
    if !ok {
        log.Fatalf("Response format error: 'items' field not found or is not a list: %v", result)
    }
    fmt.Printf(":white_check_mark: Tokens Fetched Successfully for Customer ID: %s\n", CUSTOMER_ID)
    fmt.Printf("Total Tokens Found: %d\n", len(items))
    // 7. Iterate and display key token details
    for i, item := range items {
        token, tOK := item.(map[string]interface{})
        if !tOK {
            log.Printf("Warning: Skipping malformed token at index %d", i)
            continue
        }
        tokenID := token["id"]
        cardLast4 := ""
        cardNetwork := ""
        // Token entity contains card details inside the 'card' object
        if cardData, cOK := token["card"].(map[string]interface{}); cOK {
            if last4, lOK := cardData["last4"].(string); lOK {
                cardLast4 = last4
            }
            if network, nOK := cardData["network"].(string); nOK {
                cardNetwork = network
            }
        }
        fmt.Println("\n----------------------------------------------")
        fmt.Printf("Token %d (ID: %s)\n", i+1, tokenID)
        fmt.Printf("  Card: %s ending in %s\n", cardNetwork, cardLast4)
        fmt.Printf("  Status: %s\n", token["status"])
    }
}







