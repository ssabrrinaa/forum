package models

import "github.com/gofrs/uuid"

type Category struct {
	ID   uuid.UUID
	Name string
}
