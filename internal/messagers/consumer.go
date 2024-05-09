package messagers

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ConsumerInterface interface {
	ConsumeMessage(ctx context.Context, messages chan<- proto.Message, prototype proto.Message) error
	Close() error
	exponentialBackoff(min, max, attempt int) time.Duration
}

type Consumer struct {
	reader *kafka.Reader
	logger *zap.Logger
}

type ConsumerOptions struct {
	BoostrapServers []string
	Topic           string
	GroupID         string
	Logger          *zap.Logger
}

func NewConsumer(options ConsumerOptions) (ConsumerInterface, error) {
	if len(options.BoostrapServers) == 0 {
		return nil, fmt.Errorf("brokers list must not be empty")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        options.BoostrapServers,
		Topic:          options.Topic,
		GroupID:        options.GroupID,
		MinBytes:       10e3,             // 10KB
		MaxBytes:       10e6,             // 10MB
		CommitInterval: time.Duration(0), // disable auto-commit
	})

	return &Consumer{
		reader: reader,
		logger: options.Logger,
	}, nil
}

func (c *Consumer) ConsumeMessage(ctx context.Context, messages chan<- proto.Message, prototype proto.Message) error {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Consumer stopping; context cancelled")
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("Failed to read message", zap.Error(err))
				time.Sleep(c.exponentialBackoff(100, 10000, rand.Intn(5)))
				continue
			}

			message := proto.Clone(prototype).ProtoReflect().Interface()
			if err := proto.Unmarshal(m.Value, message); err != nil {
				c.logger.Error("Failed to unmarshal protobuf message", zap.Error(err))
				continue
			}
			messages <- message

			c.reader.CommitMessages(ctx, m)
		}
	}
}

func (c *Consumer) exponentialBackoff(min, max, attempt int) time.Duration {
	minDuration := time.Duration(min) * time.Millisecond
	maxDuration := time.Duration(max) * time.Millisecond
	delay := minDuration + time.Duration(rand.Intn(int(maxDuration-minDuration+1)))
	return time.Duration(float64(delay) * math.Pow(2, float64(attempt)))
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		c.logger.Error("Failed to close Kafka reader", zap.Error(err))
		return err
	}
	c.logger.Info("Consumer connection closed")
	return nil
}
