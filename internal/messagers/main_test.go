package messagers_test

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"log"
	"os"
	"testing"
)

var kafkaBootstrapServers string

func TestMain(m *testing.M) {
	ctx := context.Background()

	kafkaContainer, err := kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/cp-kafka:7.6.1"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container after
	defer func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	host, err := kafkaContainer.Host(ctx)
	if err != nil {
		panic(err)
	}
	port, err := kafkaContainer.MappedPort(ctx, "9092")
	if err != nil {
		panic(err)
	}
	kafkaBootstrapServers = fmt.Sprintf("%s:%s", host, port.Port())

	code := m.Run()
	os.Exit(code)
}
