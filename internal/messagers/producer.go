package messagers

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ProducerInterface interface {
	Start(ctx context.Context) error
	Close() error
	CreateTopic(topic string) error
}

type Producer struct {
	writer      *kafka.Writer
	logger      *zap.Logger
	readChannel <-chan PublishInput
	connection  *kafka.Conn
}

type ProducerOptions struct {
	BoostrapServers []string
	Logger          *zap.Logger
	ReadChannel     <-chan PublishInput
}

type PublishInput struct {
	Key     []byte
	Message proto.Message
	Topic   string
}

func NewProducer(options ProducerOptions) (ProducerInterface, error) {
	if len(options.BoostrapServers) == 0 {
		return nil, fmt.Errorf("brokers list must not be empty")
	}

	connection, err := kafka.Dial("tcp", options.BoostrapServers[0])
	if err != nil {
		options.Logger.Error("Failed to connect to Kafka", zap.Error(err))
		return nil, err

	}

	writer := &kafka.Writer{
		Addr: kafka.TCP(options.BoostrapServers...),
	}

	options.Logger.Info("Kafka producer created")

	return &Producer{
		writer:      writer,
		logger:      options.Logger,
		readChannel: options.ReadChannel,
		connection:  connection,
	}, nil
}

func (p *Producer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case input, ok := <-p.readChannel:
			if !ok {
				p.logger.Info("Read channel closed, stopping producer")
				return nil
			}
			if err := p.publishMessage(ctx, input); err != nil {
				p.logger.Error("Failed to publish message", zap.Error(err))

				// TODO WE DONT NEED NO RETRY
			}
		}
	}
}

func (p *Producer) publishMessage(ctx context.Context, input PublishInput) error {
	body, err := proto.Marshal(input.Message)
	if err != nil {
		p.logger.Error("Failed to marshal protobuf message", zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: body,
		Key:   input.Key,
		Topic: input.Topic,
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

func (p *Producer) CreateTopic(topic string) error {
	err := p.connection.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		p.logger.Error("Failed to create topic", zap.Error(err))
		return err
	}
	p.logger.Info("Topic created")
	return nil
}
