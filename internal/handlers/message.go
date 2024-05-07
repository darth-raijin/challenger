package handlers

import (
	"fmt"
	"github.com/darth-raijin/challenger/internal/protos"
	"github.com/darth-raijin/challenger/internal/repositories"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

//go:generate mockery --name Message --inpackage
type Message interface {
	HandleMessage(data []byte) error
}

type MessageOptions struct {
	Logger            *zap.Logger
	MessageRepository repositories.Message
}

type message struct {
	logger            *zap.Logger
	messageRepository repositories.Message
}

func NewMessage(options MessageOptions) Message {
	return &message{
		logger:            options.Logger,
		messageRepository: options.MessageRepository,
	}
}

func (m message) HandleMessage(data []byte) error {
	var msg protos.Order
	if err := proto.Unmarshal(data, &msg); err != nil {
		m.logger.Error("Failed to unmarshal message", zap.Error(err))
		return err
	}

	fmt.Println(fmt.Sprintf("Received message: %s", msg.String()))

	return nil
}
