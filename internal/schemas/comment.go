package schemas

import (
	"forum/internal/models"

	"github.com/gofrs/uuid"
)

type CreateComment struct {
	ID      uuid.UUID
	Content string
	PostID  uuid.UUID
	UserID  uuid.UUID
}

type CreateCommentForm struct {
	TemplateCommentForm *TemplateCommentForm
	Session             *models.Session
}

type TemplateCommentForm struct {
	CommentErrors     string
	CommentDataForErr string
}
