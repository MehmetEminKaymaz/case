package main

import (
	. "caseProject/Config"
	"caseProject/Consumers"
	. "caseProject/DataModel"
	. "caseProject/RestAPI/Messages"
	"context"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm/clause"
	"log"
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	running  bool
	cancelFn context.CancelFunc
)

func startHandler(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()

	if running {
		c.WriteString("Already running")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancelFn = cancel
	running = true

	go dispatcher(ctx)

	c.WriteString("Started")
	return nil
}

func stopHandler(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()

	if !running {
		c.WriteString("Not running")
		return nil
	}

	cancelFn()
	running = false

	c.WriteString("Stopped")
	return nil
}

func dispatcher(ctx context.Context) {
	// RabbitMQ bağlantısı (bir kez)
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

	_, err = ch.QueueDeclare(
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

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Dispatcher stopped")
			return
		case <-ticker.C:
			var outboxes []MessageOutbox
			err := Database.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
				Where("dispatched = ?", false).
				Limit(2).
				Find(&outboxes).Error
			if err != nil {
				log.Println("Error fetching outbox:", err)
				continue
			}

			for _, out := range outboxes {
				err := ch.Publish(
					"",
					"message_queue",
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        []byte(out.Payload),
					},
				)
				if err != nil {
					log.Println("Failed to publish message:", err)
					continue
				}

				now := time.Now()
				if err := Database.Model(&out).Updates(map[string]interface{}{
					"dispatched":    true,
					"dispatched_at": &now,
				}).Error; err != nil {
					log.Println("Failed to update dispatched status:", err)
					continue
				}

				log.Printf("Dispatched message ID: %s\n", out.MessageId)
			}

			log.Println("Worker tick: sending messages...")
		}
	}
}

func main() {
	InitDatabaseConfig()

	go Consumers.MessageConsumer()

	app := fiber.New()

	app.Get("/start", startHandler)
	app.Get("/stop", stopHandler)

	app.Post("/message", CreateMessage)
	app.Get("/messages", GetMessages)

	log.Fatal(app.Listen(":8080"))
}
