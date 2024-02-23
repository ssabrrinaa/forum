package posthandler

import (
	"fmt"
	"forum/internal/exceptions"
	. "forum/internal/models"
	"forum/internal/schemas"
	"forum/pkg/cust_encoders"
	"forum/pkg/validator"
	"github.com/gofrs/uuid"
	"html/template"
	"net/http"
)

func (ah *PostHandler) PostCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		sessionValue := r.Context().Value("session")

		if sessionValue == nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}

		var toOk Session
		_ = toOk

		session, ok := sessionValue.(Session)
		if !ok {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		categories, err := ah.PostService.GetAllCategories()
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		}
		createPostForm := &schemas.CreatePostForm{}
		createPostForm.Session = &session
		createPostForm.Categories = categories
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				dataErr := exceptions.NewBadRequestError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			} else {
				title := r.FormValue("title")
				body := r.FormValue("body")
				categories := r.PostForm["categories"]

				titleOk, msgTitle := validator.ValidatePostTitle(title)
				bodyOk, msgBody := validator.ValidatePostBody(body)
				categoryOk, msgCategory := validator.ValidateCategoryLen(categories)

				if !titleOk || !bodyOk || !categoryOk {
					createPostForm.TemplatePostForm = &schemas.TemplatePostForm{}
					if !titleOk {
						createPostForm.TemplatePostForm.PostErrors.Title = msgTitle
					}

					if !bodyOk {
						createPostForm.TemplatePostForm.PostErrors.Body = msgBody
					}

					if !categoryOk {
						createPostForm.TemplatePostForm.PostErrors.Category = msgCategory
					}

					createPostForm.TemplatePostForm.PostDataForErr.Title = title
					createPostForm.TemplatePostForm.PostDataForErr.Body = body

				}
				if titleOk && bodyOk && categoryOk {
					post := schemas.CreatePost{
						Title:      title,
						Body:       body,
						Image:      "/dgf/dfg",
						Likes:      0,
						Dislikes:   0,
						Categories: categories,
					}

					err = ah.PostService.CreatePost(session.UserID, post)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					http.Redirect(w, r, "/post", http.StatusSeeOther)
					return
				}
			}
		}

		t, err := template.ParseFiles("ui/templates/create.html")
		if err != nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		t.Execute(w, createPostForm)
		return
	}
	dataErr := exceptions.NewStatusMethodNotAllowed()
	params := cust_encoders.EncodeParams(dataErr)
	http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
	return
}

func (ah *PostHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	userID := session.UserID
	postID := r.FormValue("post_id")
	if postID == "" {
		dataErr := exceptions.NewBadRequestError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			dataErr := exceptions.NewBadRequestError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		} else {
			values := r.URL.Query()
			like := values.Get("like")
			dislike := values.Get("dislike")
			var isLike bool
			if like != "" && dislike != "" {
				dataErr := exceptions.NewBadRequestError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			if like != "" {
				isLike = true
			}
			postIDArg, err := uuid.FromString(postID)
			if err != nil {
				dataErr := exceptions.NewBadRequestError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			getPostResponse, err := ah.PostService.GetPost(postIDArg)
			if err != nil {
				dataErr := exceptions.NewResourceNotFoundError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			postUpdate := schemas.UpdatePost{
				PostID: postIDArg,
				CreatePost: schemas.CreatePost{
					Title:    getPostResponse.PostTitle,
					Body:     getPostResponse.PostBody,
					Image:    getPostResponse.PostImage,
					Likes:    getPostResponse.Likes,
					Dislikes: getPostResponse.Dislikes,
				},
			}
			vote, err := ah.PostService.GetVote(postIDArg, userID)
			fmt.Println(vote)

			voteCreate := schemas.CreateVote{
				ShowVote: schemas.ShowVote{
					VoteID: uuid.Must(uuid.NewV4()),
					UserID: userID,
					PostID: postIDArg,
					Binary: -1,
				},
			}
			var binary int
			if err != nil {
				if isLike {
					postUpdate.Likes++
					binary = 1
				} else {
					postUpdate.Dislikes++
				}
				err := ah.PostService.UpdatePost(userID, postUpdate)
				if err != nil {
					params := cust_encoders.EncodeParams(err)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
				voteCreate.Binary = binary
				err = ah.PostService.CreateVote(voteCreate)
				if err != nil {
					params := cust_encoders.EncodeParams(err)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
			} else {
				var isDelete bool

				if isLike && vote.Binary == 0 {
					postUpdate.Likes++
					if postUpdate.Dislikes != 0 {
						postUpdate.Dislikes--
					}
					isDelete = true
					binary = 1
				} else if !isLike && vote.Binary == 1 {
					postUpdate.Dislikes++
					if postUpdate.Likes != 0 {
						postUpdate.Likes--
					}
					isDelete = true
				}
				if isDelete {
					err := ah.PostService.DeleteVote(vote.VoteID, postUpdate)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
					voteCreate.Binary = binary
					voteCreate.VoteID = uuid.Must(uuid.NewV4())
					err = ah.PostService.CreateVote(voteCreate)
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
						return
					}
				}
			}
		}

		http.Redirect(w, r, "/post/get?post_id="+postID, http.StatusSeeOther)
	}
}

func (ah *PostHandler) PostGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	getPostResponse, err := ah.PostService.GetPost(postID)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	resp := &schemas.Data{
		Session: &session,
		Post:    getPostResponse,
	}

	t, err := template.ParseFiles("ui/templates/post.html") // different html
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, resp)
}

func (ah *PostHandler) PostGetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		dataErr := exceptions.NewStatusMethodNotAllowed()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	var session *Session
	var categories []*schemas.Category
	sessionValue := r.Context().Value("session")
	if sessionValue != nil {
		sessionVal, ok := sessionValue.(Session)
		if ok {
			session = &sessionVal
		}
	}

	category := r.URL.Query().Get("category")
	likedFilter := r.URL.Query().Get("liked")

	posts, err := ah.PostService.GetPostsAll(category)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	if likedFilter != "" {
		postsResp, err := ah.PostService.GetLikedPosts(session.UserID, posts)
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		posts = postsResp
	} else {
		categoriesResp, err := ah.PostService.GetAllCategories()
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			return
		}
		categories = categoriesResp
	}
	resp := &schemas.Data{
		Session:    session,
		Posts:      posts,
		Categories: categories,
	}

	t, err := template.ParseFiles("ui/templates/home.html")
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, resp)
}

func (ah *PostHandler) GetMyPosts(w http.ResponseWriter, r *http.Request) {
	sessionValue := r.Context().Value("session")

	if sessionValue == nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	session, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?p"+params, http.StatusSeeOther)
		return
	}
	fmt.Println("sesion,", session)

	getPostAllResponse, err := ah.PostService.GetMyPosts(session.UserID)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	resp := &schemas.Data{
		Session: &session,
		Posts:   getPostAllResponse,
	}

	t, err := template.ParseFiles("ui/templates/home.html")
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, resp)
}

//comments are going to be here

//func (ah *PostHandler) CommentCreate(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		dataErr := exceptions.NewStatusMethodNotAllowed()
//		params := cust_encoders.EncodeParams(dataErr)
//		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
//		return
//	}
//
//}
