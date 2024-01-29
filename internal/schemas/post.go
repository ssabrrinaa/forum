package schemas

import (
	"time"

	"github.com/gofrs/uuid"
)

type CreatePost struct {
	Title      string
	Body       string
	Image      string
	Categories []string
}

type UpdatePost struct {
	PostID uuid.UUID
	CreatePost
}

type GetPostResponse struct {
	Username   string
	PostID     uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	PostTitle  string
	PostBody   string
	PostImage  string
	Categories []string
	// Comments
	// Likes
}
