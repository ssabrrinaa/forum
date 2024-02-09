package authhandler

import (
	"fmt"
	"forum/internal/exceptions"
	"forum/internal/schemas"
	"forum/pkg/cust_encoders"
	"html/template"
	"net/http"
)

func (ah *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/register.html")
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
			}
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			dataErr := exceptions.NewBadRequestError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			user := schemas.CreateUser{
				UpdateUser: schemas.UpdateUser{
					Username: r.FormValue("username"),
					Email:    r.FormValue("email"),
					Password: r.FormValue("password"),
				},
				PasswordConfirm: r.Form.Get("password_confirm"),
			}
			err := ah.AuthService.CreateUser(user)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
			}
		}
	} else {
		// method not allowed
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("ui/templates/signin.html")
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
			}
		}

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			dataErr := exceptions.NewBadRequestError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		} else {
			user := schemas.AuthUser{
				Email:    r.Form.Get("email"),
				Password: r.Form.Get("password"),
			}
			session, err := ah.AuthService.CreateSession(user)
			if err != nil {
				params := cust_encoders.EncodeParams(err)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
			} else {
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
			}
		}
	} else {
		// method not allowed
		fmt.Println("method not allowed")
	}
}

func (ah *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")

	err := ah.AuthService.DeleteSession(cookie.Value)
	if err != nil {
		params := cust_encoders.EncodeParams(err)
		http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
	} else {
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
}
