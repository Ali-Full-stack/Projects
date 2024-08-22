package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-go/internal/model"

	"github.com/twmb/franz-go/pkg/kgo"
)

func (k *KafkaClient) CreateOrderKafka(ctx context.Context, order model.OrderInfo) error {
	orderByte, err :=json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order info : %v",err)
	}
	record :=kgo.Record{
		Topic: "order-events",
		Key: []byte("POST"),
		Value: orderByte,
	}

	err =k.Client.ProduceSync(ctx,&record).FirstErr()
	if err != nil {
		return fmt.Errorf("failed to produce message on Create: %v",err)
	}
	return nil
}
