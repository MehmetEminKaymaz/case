package Consumers

import (
	. "caseProject/Config"
	. "caseProject/DataModel"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func MessageConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Worker received message: %s", d.Body)
			log.Println("Message sent successfully to provider")

			payload := struct {
				Id        string `json:"id"`
				Recipient string `json:"recipient"`
				Content   string `json:"content"`
			}{}
			err = json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Println("Failed to parse message payload:", err)
				continue
			}

			err = Database.Model(&Message{}).Where("id = ?", payload.Id).Update("status", "sent").Error
			if err != nil {
				log.Println("Failed to update message status:", err)
			} else {
				log.Println("Message status updated to sent")
			}
		}
	}()

	log.Println("Worker waiting for messages...")
	<-forever
}
