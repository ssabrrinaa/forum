package schemas

import (
	"forum/internal/models"
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
	Like       int
	Dislike    int

	// Comments
}

type Data struct {
	Session *models.Session
	Post    *GetPostResponse
	Posts   []*GetPostResponse
	// Comments   []*schemas.Comment
	Categories []*Category
	// Error Error
}
type PostDataForErr struct {
	Title string
	Body  string
}

type PostErrors struct {
	Title string
	Body  string
}

type TemplatePostForm struct {
	PostErrors     PostErrors
	PostDataForErr PostDataForErr
}

type CreatePostForm struct {
	TemplatePostForm *TemplatePostForm
	Session          *models.Session
	Categories       []*Category
}
