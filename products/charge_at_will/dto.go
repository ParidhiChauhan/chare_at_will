package chargeatwill

type CreateCustomerRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
}

type SetupRequest struct {
	Amount   int    `json:"amount"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
	Method   string `json:"method"`
}

type AuthorizationPaymentRequest struct {
	OrderID string `json:"order_id"`
}
