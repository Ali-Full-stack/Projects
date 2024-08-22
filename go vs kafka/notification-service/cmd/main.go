package main

import (
	"context"
	"fmt"
	"log"
	"notification-service/internal/mongodb"
	"notification-service/kafka"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	client, err := kafka.ConnectKafka(os.Getenv("kafka_url"), "order-events")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	mongoRepo, err := mongodb.NewMongoRepo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Started Consuming messages.....")
	for {
		fetches := client.PollFetches(context.Background())
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatal(errs)
		}
		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			for _, record := range ftp.Records {
				switch string(record.Key) {
				case "POST":
					mongoRepo.AddOrderIntoMongoDB(record.Value)
				}
			}
		})
	}
}
