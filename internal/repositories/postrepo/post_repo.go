package postrepo

import (
	"database/sql"
	"forum/internal/models"

	"github.com/gofrs/uuid"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

type PostRepoI interface {
	CreatePost(post models.Post) error
	UpdatePost(post models.Post) error
	GetPost(postID uuid.UUID) (models.Post, error)
}

func (ar *PostRepo) CreatePost(post models.Post) error {
	stmt := `
		INSERT INTO posts (post_id, user_id, title, body, image) 
		VALUES (?, ?, ?, ?, ?);
	`
	if _, err := ar.db.Exec(stmt, post.ID, post.UserId, post.Title, post.Body, post.Image); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) UpdatePost(post models.Post) error {
	stmt := `
		UPDATE posts 
		SET updated_at = CURRENT_TIMESTAMP, title = ?, body = ?, image = ? 
		WHERE post_id = ?;
	`
	if _, err := ar.db.Exec(stmt, post.Title, post.Body, post.Image, post.ID); err != nil {
		return err
	}
	return nil
}

func (ar *PostRepo) GetPost(postID uuid.UUID) (models.Post, error) {
	var post models.Post
	stmt := `
		SELECT
			post_id,
			created_at,
			updated_at,
			user_id,
			title,
			body,
			image,
		FROM posts
		WHERE post_id;
	`
	raw := ar.db.QueryRow(stmt, postID)

	if err := raw.Scan(&post.ID, &post.ID, &post.CreatedAt, &post.Updeted_At, &post.UserId, &post.Title, &post.Body, &post.Image); err != nil {
		return post, err
	}
	return post, nil
}
