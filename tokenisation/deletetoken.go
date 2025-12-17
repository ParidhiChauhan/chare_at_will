package main
import (
    "fmt"
    "io"
    "log"
    "net/http"
)
// !!! IMPORTANT: The API call will execute with these specific values.
// Ensure token_RrrZERHIiWybWm is a token you wish to delete in your Test account.
const (
    // Your Razorpay Key ID
    RAZORPAY_KEY_ID    = "rzp_test_Df9PS5NyT01KEo"
    // Your Razorpay Key Secret
    RAZORPAY_KEY_SECRET = "OcyZcfdwniBefW7tnprMc4mI"
    // Razorpay API Base URL
    RAZORPAY_API_BASE  = "https://api.razorpay.com/v1"
    // The specific Customer ID (cust_...)
    CUSTOMER_ID_TO_DELETE_TOKEN = "cust_GcwirOxeZwg63m"
    // The specific Token ID (token_...)
    TOKEN_ID_TO_DELETE          = "token_RrrZERHIiWybWm"
)
// The main function executes the DELETE API call
func main() {
    // Note: The conditional check that caused the program to skip is removed.
    // Execution will now proceed to the API call.
    // 1. Construct the full API URL: DELETE /v1/customers/{customer_id}/tokens/{token_id}
    deleteURL := fmt.Sprintf("%s/customers/%s/tokens/%s",
        RAZORPAY_API_BASE,
        CUSTOMER_ID_TO_DELETE_TOKEN,
        TOKEN_ID_TO_DELETE)
    // 2. Create the HTTP request (DELETE method)
    req, err := http.NewRequest("DELETE", deleteURL, nil)
    if err != nil {
        log.Fatalf("Error creating DELETE request: %v", err)
    }
    // 3. Set Basic Authentication
    req.SetBasicAuth(RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET)
    // 4. Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error sending request to Razorpay: %v", err)
    }
    defer resp.Body.Close()
    // 5. Check the status code
    // Successful deletion returns 204 No Content
    if resp.StatusCode == http.StatusNoContent {
        fmt.Println("\n:white_check_mark: Token Deletion Successful!")
        fmt.Printf("Token ID %s for Customer ID %s has been deleted.\n", TOKEN_ID_TO_DELETE, CUSTOMER_ID_TO_DELETE_TOKEN)
    } else {
        // Read body for error details (e.g., 404 if token doesn't exist, 400 Bad Request)
        body, _ := io.ReadAll(resp.Body)
        log.Fatalf(":x: Token Deletion failed. Status: %s. Response: %s", resp.Status, body)
    }
}