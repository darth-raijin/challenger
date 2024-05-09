package messagers_test

import (
	"context"
	"testing"

	"github.com/darth-raijin/challenger/internal/messagers"
	"github.com/darth-raijin/challenger/internal/protos"
	"github.com/golang/protobuf/proto" // Make sure to use the correct protobuf package depending on your generated code
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestPublishMessage tests the happy path for publishing a message.
func TestPublishMessage(t *testing.T) {
	t.Run("can produce message", func(t *testing.T) {
		ctx := context.Background()
		logger := zap.NewNop()

		expectedTopic := uuid.NewString()

		expectedMessage := &protos.Order{
			Exchange:   "some exchange",
			BuySymbol:  "eth",
			SellSymbol: "sol",
			Type:       protos.OrderType_MARKET,
			Side:       protos.OrderSide_BUY,
			Quantity:   5,
			Price:      5,
			Status:     protos.OrderStatus_PENDING,
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

		// Deserialize the message into protos.Order
		receivedMessage := &protos.Order{}
		if err := proto.Unmarshal(m.Value, receivedMessage); err != nil {
			t.Fatalf("Failed to deserialize message: %v", err)
		}

		assert.Equal(t, expectedMessage.Type, receivedMessage.Type)
		assert.Equal(t, expectedMessage.Exchange, receivedMessage.Exchange)
		assert.Equal(t, expectedMessage.BuySymbol, receivedMessage.BuySymbol)
		assert.Equal(t, expectedMessage.SellSymbol, receivedMessage.SellSymbol)
		assert.Equal(t, expectedMessage.Side, receivedMessage.Side)
		assert.Equal(t, expectedMessage.Quantity, receivedMessage.Quantity)
		assert.Equal(t, expectedMessage.Price, receivedMessage.Price)
		assert.Equal(t, expectedMessage.Status, receivedMessage.Status)
	})
}
