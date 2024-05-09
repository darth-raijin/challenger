package messagers_test

import (
	"context"
	"github.com/darth-raijin/challenger/internal/messagers"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

// TestPublishMessage tests the happy path for publishing a message.
func TestPublishMessage(t *testing.T) {
	t.Run("can produce message", func(t *testing.T) {
		ctx := context.Background()
		logger := zap.NewNop()

		expectedTopic := uuid.NewString()
		expectedMessage := &anypb.Any{
			Value: []byte(uuid.New().String()),
		}

		expectedKey := []byte(uuid.New().String())

		producerOpts := messagers.ProducerOptions{
			BoostrapServers: []string{kafkaBootstrapServers},
			Logger:          logger,
			Topic:           expectedTopic,
			Partition:       0,
		}

		producer, err := messagers.NewProducer(producerOpts)
		if err != nil {
			t.Fatalf("Failed to create producer: %v", err)
		}
		defer producer.Close()

		if err := producer.PublishMessage(ctx, expectedKey, expectedMessage); err != nil {
			t.Fatalf("Failed to publish message: %v", err)
		}

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{kafkaBootstrapServers},
			Topic:     expectedTopic,
			Partition: 0,
			MinBytes:  10e3, // 10KB
			MaxBytes:  10e6, // 10MB
		})
		defer reader.Close()

		m, err := reader.ReadMessage(ctx)
		if err != nil {
			t.Fatalf("Failed to read message: %v", err)
		}

		if !proto.Equal(expectedMessage, &anypb.Any{Value: m.Value}) {
			t.Errorf("Mismatch in message content: got %v, want %v", m.Value, expectedMessage.Value)
		}
	})

}
