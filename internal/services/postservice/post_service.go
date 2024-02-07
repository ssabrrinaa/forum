package postservice

import (
	"forum/internal/exceptions"
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
	GetPostsAll() ([]*schemas.GetPostResponse, error)
	GetMyPosts(userID uuid.UUID) ([]*schemas.GetPostResponse, error)
	GetAllCategories() ([]*schemas.Category, error)
}

func (as *PostService) CreatePost(user_id uuid.UUID, postCreate schemas.CreatePost) error {
	post := models.Post{
		ID:     uuid.Must(uuid.NewV4()),
		UserId: user_id,
		Title:  postCreate.Title,
		Body:   postCreate.Body,
		Image:  postCreate.Image,
	}
	err := as.PostRepo.CreatePost(post)
	if err != nil {
		return exceptions.NewInternalServerError()
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
		return exceptions.NewInternalServerError()
	}
	return nil
}

func (as *PostService) GetPost(post_id uuid.UUID) (*schemas.GetPostResponse, error) {
	post, err := as.PostRepo.GetPost(post_id)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	user, err := as.AuthRepo.GetUserByUserID(post.UserId)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	categories, err := as.PostRepo.GetCategoriesByPostID(post_id)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	// get comments and likes

	return &schemas.GetPostResponse{
		Username:   user.Username,
		PostID:     post.ID,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
		PostTitle:  post.Title,
		PostBody:   post.Body,
		PostImage:  post.Image,
		Categories: categories,
	}, nil
}

func (as *PostService) GetPostsAll() ([]*schemas.GetPostResponse, error) {
	var getPostsAllResponse []*schemas.GetPostResponse

	posts, err := as.PostRepo.GetPostsAll()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	for _, post := range posts {
		categories, err := as.PostRepo.GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		user, err := as.AuthRepo.GetUserByUserID(post.UserId)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		getPostsAllResponse = append(getPostsAllResponse, &schemas.GetPostResponse{
			Username:   user.Username,
			PostID:     post.ID,
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
			PostTitle:  post.Title,
			PostBody:   post.Body,
			PostImage:  post.Image,
			Categories: categories,
		})
	}

	// You may want to get comments and likes here

	return getPostsAllResponse, nil
}

func (as *PostService) GetMyPosts(userID uuid.UUID) ([]*schemas.GetPostResponse, error) {
	var getMyPostsResponse []*schemas.GetPostResponse

	posts, err := as.PostRepo.GetMyPosts(userID)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	for _, post := range posts {
		categories, err := as.PostRepo.GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		user, err := as.AuthRepo.GetUserByUserID(post.UserId)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		getMyPostsResponse = append(getMyPostsResponse, &schemas.GetPostResponse{
			Username:   user.Username,
			PostID:     post.ID,
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
			PostTitle:  post.Title,
			PostBody:   post.Body,
			PostImage:  post.Image,
			Categories: categories,
		})
	}

	// You may want to get comments and likes here

	return getMyPostsResponse, nil
}

func (as *PostService) GetAllCategories() ([]*schemas.Category, error) {
	var categoriesResp []*schemas.Category
	categories, err := as.PostRepo.GetAllCategories()
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	for _, category := range categories {
		tempCategory := &schemas.Category{
			ID:   category.ID,
			Name: category.Name,
		}
		categoriesResp = append(categoriesResp, tempCategory)
	}

	return categoriesResp, nil
}
