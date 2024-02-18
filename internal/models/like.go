package models

import (
	"github.com/gofrs/uuid"
)

type Like struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	PostID   uuid.UUID
	Positive int
	Negative int
}
