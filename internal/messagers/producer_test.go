package messagers_test

import (
	"context"
	"sync"
	"testing"

	"github.com/darth-raijin/challenger/internal/messagers"
	"github.com/darth-raijin/challenger/internal/protos"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestPublishMessage tests the happy path for publishing a message to Kafka
func TestPublishMessage(t *testing.T) {
	t.Run("can produce message", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		logger := zap.NewNop()
		messageChannel := make(chan messagers.PublishInput)
		var wg sync.WaitGroup

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
			ReadChannel:     messageChannel,
		}

		producer, err := messagers.NewProducer(producerOpts)
		if err != nil {
			t.Fatalf("Failed to create producer: %v", err)
		}
		defer producer.Close()

		err = producer.CreateTopic(expectedTopic)
		if err != nil {
			t.Fatalf("Failed to create topic: %v", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := producer.Start(ctx); err != nil {
				t.Error("Producer failed to start:", err)
			}
		}()

		// Send a message
		messageChannel <- messagers.PublishInput{
			Key:     expectedKey,
			Message: expectedMessage,
			Topic:   expectedTopic,
		}

		// Setup a Kafka reader to read the message back
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{kafkaBootstrapServers},
			Topic:     expectedTopic,
			Partition: 0,
		})
		defer reader.Close()

		close(messageChannel)

		wg.Wait()

		// Read the message from Kafka
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			t.Fatalf("Failed to read message: %v", err)
		}

		// Deserialize the message into protos.Order
		receivedMessage := &protos.Order{}
		if err := proto.Unmarshal(m.Value, receivedMessage); err != nil {
			t.Fatalf("Failed to deserialize message: %v", err)
		}

		// Assert that the received message matches the expected message
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
