package Messages

import (
	. "caseProject/Config"
	. "caseProject/DataModel"
	"github.com/gofiber/fiber/v2"
)

func GetMessages(c *fiber.Ctx) error {

	var messages []Message

	err := Database.Model(&Message{}).
		Where("status = ?", "sent").
		Select("id, content, recipient, status").
		Scan(&messages).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve messages",
		})
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}
