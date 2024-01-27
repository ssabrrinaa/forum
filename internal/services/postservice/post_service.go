package postservice

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repositories/postrepo"
	"forum/internal/schemas"

	"github.com/gofrs/uuid"
)

type PostService struct {
	PostRepo postrepo.PostRepoI
}

func NewPostService(postRepo postrepo.PostRepoI) *PostService {
	return &PostService{
		PostRepo: postRepo,
	}
}

type PostServiceI interface {
	CreatePost(user_id uuid.UUID, postCreate schemas.CreatePost) error
	UpdatePost(user_id uuid.UUID, postCreate schemas.UpdatePost) error
	GetPost(post_id uuid.UUID) (models.Post, error)
}

func (as *PostService) CreatePost(user_id uuid.UUID, postCreate schemas.CreatePost) error {
	fmt.Println("++++++++++++++++++")
	fmt.Println(postCreate.Body)

	post := models.Post{
		ID:     uuid.Must(uuid.NewV4()),
		UserId: user_id,
		Title:  postCreate.Title,
		Body:   postCreate.Body,
		Image:  postCreate.Image,
	}
	fmt.Println(post)
	err := as.PostRepo.CreatePost(post)
	if err != nil {
		fmt.Println("PostRepo", err)
		return err
	}
	return nil
}

func (as *PostService) UpdatePost(user_id uuid.UUID, postCreate schemas.UpdatePost) error {
	post := models.Post{
		ID:     postCreate.PostID,
		UserId: user_id,
		Title:  postCreate.CreatePost.Title,
		Body:   postCreate.CreatePost.Body,
		Image:  postCreate.Image,
	}

	err := as.PostRepo.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (as *PostService) GetPost(post_id uuid.UUID) (models.Post, error) {
	post, err := as.PostRepo.GetPost(post_id)
	if err != nil {
		return post, err
	}
	fmt.Println(post)
	return post, nil
}
