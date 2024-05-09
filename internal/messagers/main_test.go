package messagers_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
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

	ports, err := kafkaContainer.MappedPort(ctx, "9093")
	if err != nil {
		panic(err)
	}

	kafkaBootstrapServers = fmt.Sprintf("%s:%s", host, ports.Port())

	code := m.Run()
	os.Exit(code)
}
