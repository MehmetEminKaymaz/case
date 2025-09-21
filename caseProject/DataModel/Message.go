package DataModel

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	Content   string
	Recipient string
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
