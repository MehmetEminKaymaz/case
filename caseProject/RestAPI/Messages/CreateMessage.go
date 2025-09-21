package Messages

import (
	. "caseProject/Config"
	. "caseProject/DataModel"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CreateMessage(c *fiber.Ctx) error {

	requestPayload := struct {
		Content   string `json:"content"`
		Recipient string `json:"recipient"`
	}{}

	if err := c.BodyParser(&requestPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	transaction := Database.Begin()

	msg := Message{
		Content:   requestPayload.Content,
		Recipient: requestPayload.Recipient,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	if err := transaction.Create(&msg).Error; err != nil {
		transaction.Rollback()
		return err
	}

	payload := fmt.Sprintf(`{"id":"%s","recipient":"%s","content":"%s"}`, msg.Id.String(), msg.Recipient, msg.Content)

	outbox := MessageOutbox{
		MessageId:    msg.Id,
		Payload:      payload,
		Topic:        "sms.send",
		Dispatched:   false,
		CreatedAt:    time.Now(),
		DispatchedAt: nil,
	}

	if err := transaction.Create(&outbox).Error; err != nil {
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit().Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}
