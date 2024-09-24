package viva

// Token represents the structure of the viva token response
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type OrderResponse struct {
	OrderCode int `json:"orderCode"`
}

// OrderRequest represents the structure of the order creation request
type OrderRequest struct {
	Amount float64 `json:"amount"`
	// CustomerTrns - a short description of the items being purchased
	CustomerTrns string   `json:"customerTrns"`
	Customer     Customer `json:"customer"`
}

// Customer represents the customer information in the order request
type Customer struct {
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
	RequestLang string `json:"requestLang"`
}
