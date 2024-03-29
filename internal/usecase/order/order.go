package create_order

import (
	"context"
	"fmt"
	"log/slog"

	"outbox-transactional/internal/pkg/entity/order"
	orderRepo "outbox-transactional/internal/pkg/repository/order"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

// const cursedUser = 666

// Usecase responsible for saving request.
type Usecase struct {
	repo     orderRepo.OrderRepo
	producer sarama.SyncProducer
}

func New(orderRepo orderRepo.OrderRepo, producer sarama.SyncProducer) *Usecase {
	return &Usecase{
		repo:     orderRepo,
		producer: producer,
	}
}

// Save single order
func (uc *Usecase) Save(ctx context.Context, log slog.Logger, order *order.Order) error {
	_, err := uc.repo.Save(ctx, order)
	if err != nil {
		return errors.Wrap(err, "repo.Save")
	}

	// if order.UserID == cursedUser {
	// 	return errors.New("some err")
	// }

	// if _, _, err = uc.producer.SendMessage(&sarama.ProducerMessage{
	// 	Topic: kafka.Topic,
	// 	Value: sarama.StringEncoder(fmt.Sprintf("{order_id:%d}", orderID)),
	// }); err != nil {
	// 	return errors.Wrap(err, "producer.SendMessage")
	// }

	// tx.Commit()

	return nil
}

// Get orders by ids
func (uc *Usecase) Get(ctx context.Context, IDs []uint64) ([]order.Order, error) {
	ordersMap, err := uc.repo.Get(ctx, IDs)
	if err != nil {
		return nil, fmt.Errorf("err from orders_repository: %s", err.Error())
	}

	// count amount and discount for all orders.
	for idx, singleOrder := range ordersMap {
		for _, singleService := range singleOrder.Items {
			singleOrder.OriginalAmount += singleService.Amount
			singleOrder.DiscountedAmount += singleService.DiscountedAmount
			ordersMap[idx] = singleOrder
		}
	}

	result := make([]order.Order, 0, len(ordersMap))
	for _, ord := range ordersMap {
		result = append(result, ord)
	}

	return result, nil
}
