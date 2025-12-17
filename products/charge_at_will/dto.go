package chargeatwill

type AuthorizationRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
	Method   string `json:"method"` // upi
}
