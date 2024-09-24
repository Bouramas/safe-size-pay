package responses

// TransactionResponse - the response returned upon creting a new transcation
type TransactionResponse struct {
	// RedirectUrl - consists of the orderCode as well , the UI will use this to show a QR Code
	RedirectUrl string `json:"redirect_url"`
	Status      string `json:"status"`
	ID string `json:"id"`
}
