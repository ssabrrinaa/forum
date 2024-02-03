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
	GetPostsAll() ([]models.Post, error)
	GetMyPosts(userID uuid.UUID) ([]models.Post, error)
	GetCategoriesByPostID(postID uuid.UUID) ([]string, error)
	GetAllCategories() ([]*models.Category, error)
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
		WHERE post_id = ?;
	`
	raw := ar.db.QueryRow(stmt, postID)

	if err := raw.Scan(&post.ID, &post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Title, &post.Body, &post.Image); err != nil {
		return post, err
	}
	return post, nil
}

func (ar *PostRepo) GetPostsAll() ([]models.Post, error) {
	var posts []models.Post
	stmt := `
		SELECT
			post_id,
			created_at,
			updated_at,
			user_id,
			title,
			body,
			image
		FROM posts;
	`

	rows, err := ar.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Title, &post.Body, &post.Image); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (ar *PostRepo) GetCategoriesByPostID(postID uuid.UUID) ([]string, error) {
	var categories []string
	stmt := `
		SELECT c.name
		FROM categories c
		JOIN categories_posts_association a ON c.category_id = a.category_id
		JOIN posts p ON a.post_id = p.post_id
		WHERE p.post_id = ?;
	`

	rows, err := ar.db.Query(stmt, postID)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			return categories, err
		}
		categories = append(categories, categoryName)
	}

	if err := rows.Err(); err != nil {
		return categories, err
	}

	return categories, nil
}

func (ar *PostRepo) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	stmt := `
		SELECT name
		FROM categories 
	`

	rows, err := ar.db.Query(stmt)
	if err != nil {
		return categories, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			return categories, err
		}
		category := &models.Category{Name: categoryName}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return categories, err
	}

	return categories, nil
}

func (ar *PostRepo) GetMyPosts(userID uuid.UUID) ([]models.Post, error) {
	var posts []models.Post
	stmt := `
		SELECT
			post_id,
			created_at,
			updated_at,
			user_id,
			title,
			body,
			image
		FROM posts
		WHERE user_id = ?;
	`

	rows, err := ar.db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserId, &post.Title, &post.Body, &post.Image); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
