package rabbitmq

import (
	"encoding/json"
	"log"
	"rabbitmq-topic/internal/model"
	"rabbitmq-topic/internal/mongodb"
)

func (r *RabbitRepo) ConsumeReports(m *mongodb.MongoRepo){
	q, err := r.RabbitChannel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}
	reports :=[]string{"*.created", "*.updated", "*.deleted"}
	for _, s := range reports {
		err = r.RabbitChannel.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"report-exchange", // exchange
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("%s: %s", "Failed to bind a queue", err)
		}
	}

	msgs, err := r.RabbitChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			status :=d.Headers["type"].(string)
			switch status {
			case "created":
				var  reportDetails model.ReportDetails
				if err :=json.Unmarshal(d.Body, &reportDetails); err != nil {
					log.Println("Failed to unmarshal message body :",err)
				}
				if err :=m.AddNewReportIntoMongoDB(reportDetails); err != nil {
					log.Println(err)
				}
			case "updated":
				id :=d.Headers["id"].(string)
				var  reportDetails model.ReportDetails
				if err :=json.Unmarshal(d.Body, &reportDetails); err != nil {
					log.Println("Failed to unmarshal message body:",err)
				}
				if err :=m.UpdateExistingReportInMongoDB(reportDetails, id); err != nil {
					log.Println(err)
				}
			case "deleted":
				id :=d.Headers["id"].(string)
				if err :=m.DeleteExistingReportFromMongoDB(id); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for Reports .....")
	<-forever
}
