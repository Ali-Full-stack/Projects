package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

func ConnectKafka(kurl string, topic string, partition, replication int) (*KafkaClient, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(kurl),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed kafka Client error: %v", err)
	}
	err = client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed kafka connection: %v", err)
	}
	// admin := kadm.NewClient(client)
	// _, err = admin.CreateTopic(context.Background(), int32(partition), int16(replication), nil, topic)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create topic: %v", err)
	// }
	return &KafkaClient{Client: client}, nil
}

type KafkaClient struct {
	Client *kgo.Client
}
