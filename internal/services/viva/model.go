package viva

import (
	"math"
	"time"
)

// Token represents the structure of the viva token response
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type CreateOrderResponse struct {
	OrderCode int `json:"orderCode"`
}

// CreateOrderRequest represents the structure of the order creation request
type CreateOrderRequest struct {
	Amount int `json:"amount"`
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

type TransactionUpdate struct {
	Email                string      `json:"email"`
	BankId               string      `json:"bankId"`
	Amount               float64     `json:"amount"`
	ConversionRate       float64     `json:"conversionRate"`
	OriginalAmount       float64     `json:"originalAmount"`
	OriginalCurrencyCode string      `json:"originalCurrencyCode"`
	SourceCode           string      `json:"sourceCode"`
	Switching            bool        `json:"switching"`
	OrderCode            int64       `json:"orderCode"`
	StatusId             string      `json:"statusId"`
	FullName             string      `json:"fullName"`
	InsDate              time.Time   `json:"insDate"`
	CardNumber           string      `json:"cardNumber"`
	CurrencyCode         string      `json:"currencyCode"`
	CustomerTrns         string      `json:"customerTrns"`
	MerchantTrns         string      `json:"merchantTrns"`
	TransactionTypeId    int         `json:"transactionTypeId"`
	RecurringSupport     bool        `json:"recurringSupport"`
	TotalInstallments    int         `json:"totalInstallments"`
	CardCountryCode      string      `json:"cardCountryCode"`
	CardUniqueReference  string      `json:"cardUniqueReference"`
	CardIssuingBank      interface{} `json:"cardIssuingBank"`
	EventId              interface{} `json:"EventId"`
	CurrentInstallment   int         `json:"currentInstallment"`
	CardTypeId           int         `json:"cardTypeId"`
	DigitalWalletId      int         `json:"digitalWalletId"`
}

// AmountConversion - converts the provided value to an int
// For example €100.37 will be converted to 10037.
func AmountConversion(value float64) int {
	value = math.Round(value * 100)
	return int(value)
}
