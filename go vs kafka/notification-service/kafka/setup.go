package kafka

import (
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

func ConnectKafka(kurl string, topic string) (*kgo.Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(kurl),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup("my-group"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed consumer client:%v", err)
	}
	return client, nil
}
