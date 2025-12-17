package chargeatwill



type AuthorizationPayload struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Method   string `json:"method"`
}
