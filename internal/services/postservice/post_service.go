package postservice

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repositories/authrepo"
	"forum/internal/repositories/postrepo"
	"forum/internal/schemas"

	"github.com/gofrs/uuid"
)

type PostService struct {
	PostRepo postrepo.PostRepoI
	AuthRepo authrepo.AuthRepoI
}

func NewPostService(postRepo postrepo.PostRepoI, authrepo authrepo.AuthRepoI) *PostService {
	return &PostService{
		PostRepo: postRepo,
		AuthRepo: authrepo,
	}
}

type PostServiceI interface {
	CreatePost(user_id uuid.UUID, postCreate schemas.CreatePost) error
	UpdatePost(user_id uuid.UUID, postCreate schemas.UpdatePost) error
	GetPost(post_id uuid.UUID) (*schemas.GetPostResponse, error)
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

func (as *PostService) GetPost(post_id uuid.UUID) (*schemas.GetPostResponse, error) {
	// var getPostResponce schemas.GetPostResponse
	post, err := as.PostRepo.GetPost(post_id)
	if err != nil {
		return nil, err
	}

	user, err := as.AuthRepo.GetUserByUserID(post.UserId)
	if err != nil {
		return nil, err
	}

	categories, err := as.PostRepo.GetCategoriesByPostID(post_id)
	if err != nil {
		return nil, err
	}

	// get comments and likes

	return &schemas.GetPostResponse{
		Username:   user.Username,
		PostID:     post.ID,
		CreatedAt:  post.CreatedAt,
		UpdetedAt:  post.UpdeatedAt,
		PostTitle:  post.Title,
		PostBody:   post.Body,
		PostImage:  post.Image,
		Categories: categories,
	}, nil
}
