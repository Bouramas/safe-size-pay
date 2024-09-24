package viva

import (
	"os"
	"testing"
)

func Test_client_GetToken(t *testing.T) {

	os.Setenv("VIVA_USERNAME", "qv90226c59cdr7kx9c4zs770ty4cf346kfv194634ek59.apps.vivapayments.com")
	os.Setenv("VIVA_PASSWORD", "Guww0WPVEbadtwgmjrSzF8d4mve16E")
	os.Setenv("VIVA_BASE_ACCOUNTS_URL", "https://demo-accounts.vivapayments.com")
	os.Setenv("VIVA_BASE_API_URL", "https://demo-api.vivapayments.com")

	vc := NewClient()
	orderRequest := &OrderRequest{
		Amount:       2.2,
		CustomerTrns: "whatever",
		Customer: Customer{
			"random@random.com",
			"Giannis Bouramas",
			"en-US",
		},
	}
	got, err := vc.CreateOrder(orderRequest)
	if err != nil {
		println(err.Error())
	}

	println("ORDER CODE: ", got.OrderCode)

	// println(got.AccessToken)
	// println(got.Scope)
	// println(got.TokenType)
	// println(got.ExpiresIn)

}
