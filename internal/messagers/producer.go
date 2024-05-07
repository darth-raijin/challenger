package messagers

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"log"
)

// ProducerInterface defines the methods a producer should have.
type ProducerInterface interface {
	PublishMessage(ctx context.Context, key []byte, message proto.Message) error
	Close() error
}

// Producer represents a message producer for Kafka.
type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
	conn   *kafka.Conn
}

// ProducerOptions holds the configuration options for setting up the Producer.
type ProducerOptions struct {
	Brokers []string
	Logger  *zap.Logger
	Topic   string
}

// NewProducer creates and returns a new Producer based on the provided options.
func NewProducer(opts ProducerOptions) (ProducerInterface, error) {
	if len(opts.Brokers) == 0 {
		return nil, fmt.Errorf("brokers list must not be empty")
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9094", opts.Topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9094"),
		Topic: opts.Topic,
	}

	opts.Logger.Info("Kafka producer created")

	return &Producer{
		writer: writer,
		logger: opts.Logger,
		conn:   conn,
	}, nil
}

// PublishMessage sends a message to the specified Kafka topic with the given key.
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

// Close closes the connection to Kafka.
func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		p.logger.Error("Failed to close Kafka writer", zap.Error(err))
		return err
	}
	p.logger.Info("Producer connection closed")
	return nil
}
