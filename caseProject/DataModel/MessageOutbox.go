package DataModel

import (
	"github.com/google/uuid"
	"time"
)

type MessageOutbox struct {
	Id           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	MessageId    uuid.UUID `gorm:"type:uuid;not null;index"`
	Payload      string    `gorm:"type:jsonb;not null"`
	Topic        string    `gorm:"type:varchar(100);not null"`
	Dispatched   bool      `gorm:"default:false"`
	CreatedAt    time.Time
	DispatchedAt *time.Time
}
