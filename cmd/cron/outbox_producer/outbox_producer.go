package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"outbox-transactional/cmd/config"
	"outbox-transactional/internal/cron/outbox_producer"
	"outbox-transactional/internal/kafka"
	"outbox-transactional/internal/pkg/repository/order"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Get logger interface.
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
		return fmt.Errorf(" 'conf' flag required")
	}

	cfg, err := config.ParseConfig(confString)
	if err != nil {
		return errors.Wrap(err, "config.Parse")
	}

	log.Info("cfg:", slog.Any("cfg", cfg))
	log.Info("Starting the service...")

	ctx := context.Background()

	producer := kafka.NewProducer(cfg.KafkaPort)

	pool, err := pgxpool.Connect(context.Background(), cfg.DbConnString)
	if err != nil {
		return fmt.Errorf("can't create pg pool: %s", err.Error())
	}

	outboxProducer := outbox_producer.New(producer, pool, order.OutboxTable, log)

	return errors.Wrap(outboxProducer.ProduceMessages(ctx), "outboxProducer.ProduceMessages")
}
