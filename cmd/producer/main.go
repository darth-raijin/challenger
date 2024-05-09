package main

import (
	"context"
	"fmt"
	"github.com/darth-raijin/challenger/internal/config"
	"github.com/darth-raijin/challenger/internal/messagers"
	"github.com/darth-raijin/challenger/internal/protos"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, err := zap.NewProduction()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Error("error loading cfg", zap.Error(err))
		os.Exit(1)
	}

	producer, err := wireDependencies(cfg, logger, err)
	if err != nil {
		logger.Error("error wiring dependencies", zap.Error(err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for i := 0; i < cfg.KafkaProducer.Workers; i++ {
		go func() {
			for {
				inputMessage := protos.Order{
					Exchange:   "",
					BuySymbol:  "",
					SellSymbol: "",
					Type:       0,
					Side:       0,
					Quantity:   0,
					Price:      0,
					Status:     0,
				}

				key := []byte(fmt.Sprintf("%v.%v.%v", inputMessage.Exchange, inputMessage.Quantity, inputMessage.Price))
				if err := producer.PublishMessage(ctx, key, &inputMessage); err != nil {
					logger.Error("error publishing message", zap.Error(err))
					stop()
				}
			}
		}()
	}

	<-ctx.Done()
	logger.Info("shutting down gracefully")

}

func wireDependencies(config config.Config, logger *zap.Logger, err error) (messagers.ProducerInterface, error) {
	return messagers.NewProducer(messagers.ProducerOptions{
		BoostrapServers: config.Kafka.BootstrapServers,
		Logger:          logger,
		Topic:           "orders.sol.eth",
		Partition:       0,
	})
}
