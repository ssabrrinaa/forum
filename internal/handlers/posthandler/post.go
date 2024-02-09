package posthandler

import (
	"fmt"
	"forum/internal/constants"
	"forum/internal/exceptions"
	. "forum/internal/models"
	"forum/internal/schemas"
	"forum/pkg/cust_encoders"
	"forum/pkg/validator"
	"github.com/gofrs/uuid"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (ah *PostHandler) PostCreate(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodGet {
		categories, err := ah.PostService.GetAllCategories()
		if err != nil {
			params := cust_encoders.EncodeParams(err)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			t, err := template.ParseFiles("ui/templates/create.html") // different html
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			resp := &schemas.Data{
				Session:    &session,
				Categories: categories,
			}
			err = t.Execute(w, resp)
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(constants.MaxFileSize); err != nil {
			dataErr := exceptions.NewBadRequestError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			// file, header, err := r.FormFile("file")
			// fmt.Println(err)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Println("filename", header.Filename)

			// defer file.Close()

			// ext := filepath.Ext(header.Filename)
			// if !constants.IsAllowedFileExtension(ext) {
			// 	log.Println("Error: File extension not allowed.")
			// 	http.Error(w, "Unsupported file type", http.StatusBadRequest)
			// 	return
			// }
			// fmt.Println(1)

			// imageFilename := uuid.Must(uuid.NewV4()).String() + ext
			// uploadDir := constants.UploadDir
			// if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			// 	os.Mkdir(uploadDir, os.ModePerm)
			// }
			// filePath := filepath.Join(uploadDir, imageFilename)
			// newFile, err := os.Create(filePath)
			// if err != nil {
			// 	log.Fatal(err)
			// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// 	return
			// }
			// defer newFile.Close()
			// fmt.Println(2)
			// _, err = io.Copy(newFile, file)
			// if err != nil {
			// 	log.Fatal(err)
			// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// 	return
			// }

			// categories, ok := r.Form["category"]
			// if !ok{
			// 	log.Fatal("error")
			// }

			post := schemas.CreatePost{
				Title: r.FormValue("title"),
				Body:  r.FormValue("body"),
				Image: "/dgf/dfg",
				// Categories: categories,
			}

			err := validator.ValidateCreatePostInput(post)
			if err != nil {
				dataErr := exceptions.NewValidationError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			err = ah.PostService.CreatePost(session.UserID, post)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			// needs to be redirected to the right page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		// method not allowed
		fmt.Println("method not allowed")
	}
}

func (ah *PostHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/signin.html") // different html
		if err != nil {
			dataErr := exceptions.NewInternalServerError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			err := t.Execute(w, nil)
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(constants.MaxFileSize); err != nil {
			dataErr := exceptions.NewBadRequestError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			file, header, err := r.FormFile("image")

			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					dataErr := exceptions.NewInternalServerError()
					params := cust_encoders.EncodeParams(dataErr)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
			}(file)

			if err != nil {
				dataErr := exceptions.NewValidationError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			ext := filepath.Ext(header.Filename)
			if !constants.IsAllowedFileExtension(ext) {
				dataErr := exceptions.NewValidationError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			imageFilename := uuid.Must(uuid.NewV4()).String() + ext
			uploadDir := constants.UploadDir
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				err := os.Mkdir(uploadDir, os.ModePerm)
				if err != nil {
					dataErr := exceptions.NewInternalServerError()
					params := cust_encoders.EncodeParams(dataErr)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
					return
				}
			}
			filePath := filepath.Join(uploadDir, imageFilename)
			newFile, err := os.Create(filePath)
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, file)
			if err != nil {
				dataErr := exceptions.NewInternalServerError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			// categories, ok := r.Form["category"]
			// if !ok{
			// 	log.Fatal("error")
			// }
			postID, err := uuid.FromString(r.FormValue("post_id"))
			if err != nil {
				dataErr := exceptions.NewValidationError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			post := schemas.UpdatePost{
				PostID: postID,
				CreatePost: schemas.CreatePost{
					Title: r.FormValue("title"),
					Body:  r.FormValue("body"),
					Image: filePath,
					// Categories: categories,
				},
			}

			err = validator.ValidateUpdatePostInput(post)
			if err != nil {
				dataErr := exceptions.NewValidationError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			err = ah.PostService.UpdatePost(session.UserID, post)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}
			// needs to be changed, because / is busy for errors.
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	} else {
		// method not allowed
		fmt.Println("method not allowed")
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

	_, ok := sessionValue.(Session)
	if !ok {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Path[len("/post/"):]

	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	getPostResponce, err := ah.PostService.GetPost(postID)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("ui/templates/signin.html") // different html
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, getPostResponce)
}

func (ah *PostHandler) PostGetAll(w http.ResponseWriter, r *http.Request) {
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

	getPostAllResponce, err := ah.PostService.GetPostsAll()
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	categories, err := ah.PostService.GetAllCategories()
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	resp := &schemas.Data{
		Session:    &session,
		Posts:      getPostAllResponce,
		Categories: categories,
	}

	t, err := template.ParseFiles("ui/templates/home.html") // different html
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

	getPostAllResponce, err := ah.PostService.GetMyPosts(session.UserID)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("ui/templates/signin.html") // different html
	if err != nil {
		dataErr := exceptions.NewInternalServerError()
		params := cust_encoders.EncodeParams(dataErr)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		return
	}
	t.Execute(w, getPostAllResponce)
}
