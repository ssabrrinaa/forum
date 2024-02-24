package postservice

import (
	"errors"
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
	CreatePost(userID uuid.UUID, postCreate schemas.CreatePost) error
	UpdatePost(userID uuid.UUID, postUpdate schemas.UpdatePost) error
	GetPost(postID uuid.UUID) (*schemas.GetPostResponse, error)
	GetPostsAll(category string) ([]*schemas.GetPostResponse, error)
	GetMyPosts(userID uuid.UUID) ([]*schemas.GetPostResponse, error)

	GetAllCategories() ([]*schemas.Category, error)

	GetVote(postID uuid.UUID, userID uuid.UUID) (schemas.ShowVote, error)
	DeleteVote(voteID uuid.UUID, postUpdate schemas.UpdatePost) error
	CreateVote(voteCreate schemas.CreateVote) error
	GetLikedPosts(userID uuid.UUID, posts []*schemas.GetPostResponse) ([]*schemas.GetPostResponse, error)

	CreateComment(commentCreate schemas.CreateComment) error
}

func (as *PostService) CreatePost(userID uuid.UUID, postCreate schemas.CreatePost) error {
	postID := uuid.Must(uuid.NewV4())

	post := models.Post{
		ID:       postID,
		UserId:   userID,
		Title:    postCreate.Title,
		Body:     postCreate.Body,
		Likes:    postCreate.Likes,
		Dislikes: postCreate.Dislikes,
		Image:    postCreate.Image,
	}
	err := as.PostRepo.CreatePost(post)
	if err != nil {
		return exceptions.NewInternalServerError()
	}

	for _, category := range postCreate.Categories {
		err := as.PostRepo.CreatePostCategories(models.CreateCategoryPost{
			ID:           uuid.Must(uuid.NewV4()),
			CategoryName: category,
			PostID:       postID,
		})
		if err != nil {
			return exceptions.NewInternalServerError()
		}

	}

	return nil
}

func (as *PostService) UpdatePost(userID uuid.UUID, postUpdate schemas.UpdatePost) error {
	post := models.Post{
		ID:       postUpdate.PostID,
		UserId:   userID,
		Title:    postUpdate.CreatePost.Title,
		Body:     postUpdate.CreatePost.Body,
		Likes:    postUpdate.Likes,
		Dislikes: postUpdate.Dislikes,
		Image:    postUpdate.Image,
	}

	err := as.PostRepo.UpdatePost(post)
	if err != nil {
		return exceptions.NewInternalServerError()
	}
	return nil
}

func (as *PostService) GetPost(postID uuid.UUID) (*schemas.GetPostResponse, error) {
	post, err := as.PostRepo.GetPost(postID)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	user, err := as.AuthRepo.GetUserByUserID(post.UserId)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	categories, err := as.PostRepo.GetCategoriesByPostID(postID)
	if err != nil {
		return nil, exceptions.NewInternalServerError()
	}

	comments, err := as.PostRepo.GetCommentsByPostID(postID)
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
		Likes:      post.Likes,
		Dislikes:   post.Dislikes,
		PostImage:  post.Image,
		Categories: categories,
		Comments:   comments,
	}, nil
}

func (as *PostService) GetVote(postID uuid.UUID, userID uuid.UUID) (schemas.ShowVote, error) {
	voteResponse := schemas.ShowVote{}
	vote, err := as.PostRepo.GetVoteOfPost(postID, userID)
	if err != nil {
		return voteResponse, errors.New("vote is not found")
	}

	voteResponse.VoteID = vote.ID
	voteResponse.UserID = vote.UserID
	voteResponse.PostID = vote.PostID
	voteResponse.Binary = vote.Binary
	return voteResponse, nil
}

func (as *PostService) DeleteVote(voteID uuid.UUID, postUpdate schemas.UpdatePost) error {
	err := as.PostRepo.DeleteVoteOfPost(voteID)
	if err != nil {
		return err
	}
	post := models.Post{
		ID:       postUpdate.PostID,
		Title:    postUpdate.Title,
		Body:     postUpdate.Body,
		Image:    postUpdate.Image,
		Likes:    postUpdate.Likes,
		Dislikes: postUpdate.Dislikes,
	}
	err = as.PostRepo.UpdatePost(post)
	return err
}

func (as *PostService) CreateVote(voteCreate schemas.CreateVote) error {
	vote := models.Vote{
		ID:     voteCreate.VoteID,
		UserID: voteCreate.UserID,
		PostID: voteCreate.PostID,
		Binary: voteCreate.Binary,
	}
	err := as.PostRepo.CreateVote(vote)
	if err != nil {
		return err
	}
	return nil
}

func (as *PostService) GetPostsAll(category string) ([]*schemas.GetPostResponse, error) {
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

		if category != "" {
			categoryFound := false
			for _, cat := range categories {
				if cat == category {
					categoryFound = true
					break
				}
			}
			if !categoryFound {
				continue
			}
		}

		user, err := as.AuthRepo.GetUserByUserID(post.UserId)
		if err != nil {
			return nil, exceptions.NewInternalServerError()
		}

		response := &schemas.GetPostResponse{
			Username:   user.Username,
			PostID:     post.ID,
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
			PostTitle:  post.Title,
			PostBody:   post.Body,
			Likes:      post.Likes,
			Dislikes:   post.Dislikes,
			PostImage:  post.Image,
			Categories: categories,
		}

		getPostsAllResponse = append(getPostsAllResponse, response)
	}

	// You may want to get comments and likes here

	return getPostsAllResponse, nil
}

func (as *PostService) GetLikedPosts(userID uuid.UUID, posts []*schemas.GetPostResponse) ([]*schemas.GetPostResponse, error) {
	votes, err := as.PostRepo.GetVotes()
	if err != nil {
		return []*schemas.GetPostResponse{}, err
	}
	var postsResp []*schemas.GetPostResponse
	for _, post := range posts {
		for _, vote := range votes {
			if vote.UserID == userID && post.PostID == vote.PostID && vote.Binary == 1 {
				postsResp = append(postsResp, post)
			}
		}
	}
	return postsResp, nil
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

func (as *PostService) CreateComment(commentCreate schemas.CreateComment) error {
	commentID := uuid.Must(uuid.NewV4())

	comment := models.Comment{
		ID:          commentID,
		UserID:      commentCreate.UserID,
		Description: commentCreate.Content,
		PostID:      commentCreate.PostID,
	}
	err := as.PostRepo.CreateComment(comment)
	if err != nil {
		return exceptions.NewInternalServerError()
	}

	// for _, category := range postCreate.Categories {
	// 	err := as.PostRepo.CreatePostCategories(models.CreateCategoryPost{
	// 		ID:           uuid.Must(uuid.NewV4()),
	// 		CategoryName: category,
	// 		PostID:       postID,
	// 	})
	// 	if err != nil {
	// 		return exceptions.NewInternalServerError()
	// 	}

	// }

	return nil
}
