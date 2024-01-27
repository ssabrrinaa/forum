package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	Updeted_At time.Time
	UserId     uuid.UUID
	Title      string
	Body       string
	Image      string
	Categories []*Category
}
