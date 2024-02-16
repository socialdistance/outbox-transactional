package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"outbox-transactional/cmd/config"
	"outbox-transactional/internal/kafka"

	orderRepo "outbox-transactional/internal/pkg/repository/order"
	orderUCase "outbox-transactional/internal/usecase/order"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	echoRoute   = "/echo"
	orderRoute  = "/order"
	ordersRoute = "/orders"
)

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	if err := mainNoExit(log); err != nil {
		log.Info("fatal error", slog.Any("error", err))
	}
}

func mainNoExit(log *slog.Logger) error {
	// get application cfg
	confFlag := flag.String("conf", "", "cfg yaml file")
	flag.Parse()

	confString := *confFlag
	if confString == "" {
		return fmt.Errorf("conf flag required")
	}

	cfg, err := config.ParseConfig(confString)
	if err != nil {
		return errors.Wrap(err, "config.Parse")
	}

	log.Info("cfg:", slog.Any("cfg", cfg))
	log.Info("Starting the service...")

	// ctx := context.Background()

	r := mux.NewRouter()

	// echo
	// r.HandleFunc(echoRoute, echo_handler.Handler("Your message: ").ServeHTTP).Methods("GET")

	pool, err := pgxpool.Connect(context.Background(), cfg.DbConnString)
	if err != nil {
		return fmt.Errorf("can't create pg pool: %s", err.Error())
	}

	repo := orderRepo.New(pool)
	orderUseCase := orderUCase.New(repo, kafka.NewProducer(cfg.KafkaPort))

	// createOrderHandleFunc := create_order_handler.New(orderUseCase, log).Create(ctx).ServeHTTP
	// // create order
	// r.HandleFunc(orderRoute, createOrderHandleFunc).Methods("POST")

	// getOrderHandlerFunc := get_orders_handler.New(orderUseCase, log).Get(ctx).ServeHTTP
	// // get orders
	// r.HandleFunc(ordersRoute, getOrderHandlerFunc).Methods("GET")

	// if err != nil {
	// 	return fmt.Errorf("can't init router: %s", err.Error())
	// }

	log.Info("The service is ready to listen and serve.")

	return errors.Wrap(http.ListenAndServe(
		cfg.AppPort,
		r,
	), "http.ListenAndServe")
}
