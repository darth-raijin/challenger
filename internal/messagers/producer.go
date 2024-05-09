package messagers

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ProducerInterface interface {
	PublishMessage(ctx context.Context, key []byte, message proto.Message) error
	Close() error
}

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

type ProducerOptions struct {
	BoostrapServers []string
	Logger          *zap.Logger
	Topic           string
	Partition       int
}

type ProducerConfig struct {
	Partition int `mapstructure:"KAFKA_PRODUCER_PARTITION"`
	Workers   int `mapstructure:"KAFKA_PRODUCER_WORKERS"`
}

func NewProducer(opts ProducerOptions) (ProducerInterface, error) {
	if len(opts.BoostrapServers) == 0 {
		return nil, fmt.Errorf("brokers list must not be empty")
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", opts.BoostrapServers[0], opts.Topic, opts.Partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	writer := &kafka.Writer{
		Addr:  conn.RemoteAddr(),
		Topic: opts.Topic,
	}

	opts.Logger.Info("Kafka producer created")

	return &Producer{
		writer: writer,
		logger: opts.Logger,
	}, nil
}

func (p *Producer) PublishMessage(ctx context.Context, key []byte, message proto.Message) error {
	body, err := proto.Marshal(message)
	if err != nil {
		p.logger.Error("Failed to marshal protobuf message", zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: body,
		Key:   key,
	})
	if err != nil {
		p.logger.Error("Failed to publish message", zap.Error(err))
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		p.logger.Error("Failed to close Kafka writer", zap.Error(err))
		return err
	}
	p.logger.Info("Producer connection closed")
	return nil
}
