package resources

import (
	"time"
)

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	OrderID     *int      `json:"order_id,omitempty"`
	OrderMsg    string    `json:"order_msg,omitempty"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	OrderStatus string    `json:"order_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
