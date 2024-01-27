package authhandler

import (
	"fmt"
	"forum/internal/schemas"
	"html/template"
	"log"
	"net/http"
)

func (ah *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// user := &models.User{}
	fmt.Println("register")
	if r.Method == http.MethodGet {
		fmt.Println("get method")
		t, err := template.ParseFiles("ui/templates/register.html")
		if err != nil {
			log.Fatal(err) // handle the errors properly
		}
		t.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		fmt.Println("post method")
		if err := r.ParseForm(); err != nil {
			log.Fatal(err) // handle the errors properly
		}

		user := schemas.CreateUser{
			UpdateUser: schemas.UpdateUser{
				Username: r.FormValue("username"),
				Email:    r.FormValue("email"),
				Password: r.FormValue("password"),
			},
			PasswordConfirm: r.Form.Get("password_confirm"),
		}
		fmt.Println(user)
		// как чекать если юзер сущетсвует?
		err := ah.AuthService.CreateUser(user)
		if err != nil {
			log.Fatal(err) // handle the errors properly
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)

	} else {
		// method not allowed
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/signin.html")
		if err != nil {
			log.Fatal(err) // handle the errors properly
		}
		t.Execute(w, nil)
		return

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err) // handle the errors properly
		}
		user := schemas.AuthUser{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}
		session, err := ah.AuthService.CreateSession(user)
		if err != nil {
			log.Fatal(err) // handle the errors properly
		}
		cookie := &http.Cookie{
			Name:     "session",
			Value:    session.Token,
			Path:     "/",
			Expires:  session.ExpireTime,
			HttpOnly: true,
			MaxAge:   7200,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		// method not allowed
		fmt.Println("method not allowed")
	}
}

func (ah *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// Get the session cookie
	cookie, _ := r.Cookie("session")

	// Delete the session from the server side (or invalidate the session)
	err := ah.AuthService.DeleteSession(cookie.Value)
	if err != nil {
		log.Fatal(err) // Handle the error properly, e.g., return an error response
	}
	expiredCookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
