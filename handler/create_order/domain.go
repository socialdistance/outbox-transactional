package order_handler

import (
	"log/slog"
	create_order "outbox-transactional/internal/usecase/order"
)

// Handler creates orders
type Handler struct {
	uCase *create_order.Usecase
	log   slog.Logger
}

// New gives Handler.
func New(
	uCase *create_order.Usecase,
	log *slog.Logger,
) *Handler {
	return &Handler{
		uCase: uCase,
		log:   *log,
	}
}

// OrderIn is dto for http req.
type OrderIn struct {
	UserID      uint64 `json:"user_id"`
	PaymentType string `json:"payment_type"`
	Items       []Item `json:"items"`
}

type Item struct {
	ID       uint64 `json:"id"`
	Amount   uint64 `json:"amount"`
	Discount uint64 `json:"discount"`
}

var paymentTypes = map[string]PaymentType{
	"card":   Card,
	"wallet": Wallet,
}

type PaymentType uint8

const (
	UndefinedType PaymentType = iota
	Card
	Wallet
)
