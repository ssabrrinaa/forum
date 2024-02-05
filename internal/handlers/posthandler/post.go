package posthandler

import (
	"fmt"
	"forum/internal/constants"
	"forum/internal/schemas"
	"forum/pkg/validator"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
)

func (ah *PostHandler) PostCreate(w http.ResponseWriter, r *http.Request) {
	userValue := r.Context().Value("user_id")

	if userValue == nil {
		log.Fatal("User ID not found in context")
	}

	userID, ok := userValue.(uuid.UUID)
	if !ok {
		w.Header().Set("Error", "400")
		http.Redirect(w, r, "/errors", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/postcreate.html")
		if err != nil {
			w.Header().Set("Error", "500")
			http.Redirect(w, r, "/errors", http.StatusSeeOther)
		} else {
			err := t.Execute(w, nil)
			if err != nil {
				return
			}
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(constants.MaxFileSize); err != nil {
			w.Header().Set("Error", "400")
			http.Redirect(w, r, "/errors", http.StatusSeeOther)
			return
		}
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
			w.Header().Set("Error", strings.Split(err.Error(), " ")[0])
			http.Redirect(w, r, "/errors", http.StatusSeeOther)
			return
		}

		err = ah.PostService.CreatePost(userID, post)
		if err != nil {
			w.Header().Set("Error", strings.Split(err.Error(), " ")[0])
			http.Redirect(w, r, "/errors", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// method not allowed
		fmt.Println("method not allowed")
	}
}

func (ah *PostHandler) PostUpdate(w http.ResponseWriter, r *http.Request) {
	userValue := r.Context().Value("user_id")

	if userValue == nil {
		log.Fatal("User ID not found in context")
	}

	userID, ok := userValue.(uuid.UUID)
	if !ok {
		w.Header().Set("Error", "400")
		http.Redirect(w, r, "/errors", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/register.html")
		if err != nil {
			w.Header().Set("Error", "500")
			http.Redirect(w, r, "/errors", http.StatusSeeOther)
		} else {
			err := t.Execute(w, nil)
			if err != nil {
				return
			}
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(constants.MaxFileSize); err != nil {
			w.Header().Set("Error", "400")
			http.Redirect(w, r, "/errors", http.StatusSeeOther)
			return
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		ext := filepath.Ext(header.Filename)
		if !constants.IsAllowedFileExtension(ext) {
			log.Println("Error: File extension not allowed.")
			http.Error(w, "Unsupported file type", http.StatusBadRequest)
			return
		}

		imageFilename := uuid.Must(uuid.NewV4()).String() + ext
		uploadDir := constants.UploadDir
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}
		filePath := filepath.Join(uploadDir, imageFilename)
		newFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// categories, ok := r.Form["category"]
		// if !ok{
		// 	log.Fatal("error")
		// }
		postID, err := uuid.FromString(r.FormValue("post_id"))
		if err != nil {
			log.Fatal("invalid postID")
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
			log.Fatal(err) // handle the errors properly
		}

		err = ah.PostService.UpdatePost(userID, post)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// method not allowed
		fmt.Println("method not allowed")
	}
}

func (ah *PostHandler) PostGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Fatal("")
	}
	userValue := r.Context().Value("user_id")

	if userValue == nil {
		log.Fatal("User ID not found in context")
	}

	userID, ok := userValue.(uuid.UUID)
	if !ok {
		log.Fatal("Invalid user ID type in context")
	}
	postIDStr := r.URL.Path[len("/post/"):]
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		log.Fatal("invalid postID")
	}
	fmt.Println(userID)
	// validate post ID

	post, err := ah.PostService.GetPost(postID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(post)

	// http.Redirect(w, r, "/", http.StatusSeeOther)

	t, err := template.ParseFiles("ui/templates/signin.html") // different html
	if err != nil {
		log.Fatal(err) // handle the errors properly
	}
	t.Execute(w, nil)
}
