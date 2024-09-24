package requests

import (
	"testing"
)

func TestTransactionRequest_Validate(t *testing.T) {

	tests := []struct {
		name    string
		tr      *TransactionRequest
		wantErr bool
	}{
		{"success", &TransactionRequest{Amount: 10.2, Description: "random desc"}, false},
		{"description required", &TransactionRequest{Amount: 10.2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
