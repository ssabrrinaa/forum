package models

import "github.com/gofrs/uuid"

type Comment struct {
	ID          uuid.UUID
	Description string
	PostId      uuid.UUID
	UserId      uuid.UUID
}
