package order_handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	order_ucase "outbox-transactional/internal/usecase/order"
)

// Requst validation errors.
var ErrEmptyOrderIDs = errors.New("no order ids passed")

// Handler creates orders
type Handler struct {
	uCase *order_ucase.Usecase
	log   slog.Logger
}

// New gives Handler.
func New(uCase *order_ucase.Usecase, log *slog.Logger) *Handler {
	return &Handler{
		uCase: uCase,
		log:   *log,
	}
}

// GetOrdersIn is dto for http req.
type GetOrdersIn struct {
	IDs []uint64 `json:"ids"`
}

// validates request.
func (h Handler) validateReq(in *GetOrdersIn) error {
	if len(in.IDs) == 0 {
		return ErrEmptyOrderIDs
	}

	return nil
}

// Create responsible for saving new order.
func (h Handler) Get(ctx context.Context) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// prepare dto to parse request
		in := &GetOrdersIn{}
		// parse req body to dto
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			h.log.Error("can't parse req: %s", err)
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)

			return
		}

		// check that request valid
		err = h.validateReq(in)
		if err != nil {
			h.log.Error("bad req: %v: %s", in, err)
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)

			return
		}

		orders, err := h.uCase.Get(ctx, in.IDs)
		if err != nil {
			h.log.Error("can't get orders: %s", err)
			http.Error(w, "can't get orders: "+err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(orders); err != nil {
			h.log.Error("Encode: %v", err)
		}
	}

	return http.HandlerFunc(fn)
}
