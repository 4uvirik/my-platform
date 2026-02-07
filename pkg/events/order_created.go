package events

import "time"

type OrderCreated struct {
	OrderID string    `json:"order_id"`
	UserID  string    `json:"user_id"`
	Amount  int64     `json:"amount"`
	At      time.Time `json:"at"`
}
