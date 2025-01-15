package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"

	"DeliveryTimePrediction/internal/domain"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type storage interface {
	InsertResult(ctx context.Context, id string, deliveryTime float64) error
}

type ConsumerGroupHandler struct {
	ready   chan bool
	storage storage
}

func NewConsumerGroupHandler(storage storage) sarama.ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		ready:   make(chan bool),
		storage: storage,
	}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			err := h.processMessage(message)
			if err != nil {
				log.Printf("Error processing message with err %+v", err)
			}

			session.MarkMessage(message, "")
			session.Commit()
		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	var res domain.Result
	err := json.Unmarshal(message.Value, &res)
	if err != nil {
		log.Println("[consumer-group-handler] unmarshal message error:", err.Error())
		return err
	}
	err = h.storage.InsertResult(context.Background(), res.ID, res.DeliveryTime)
	if err != nil {
		log.Println("[consumer-group-handler] notify message error:", err.Error())
		return err
	}
	return nil
}
