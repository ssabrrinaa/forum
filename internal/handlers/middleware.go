package handler

import (
	"context"
	"fmt"
	"net/http"

	"forum/internal/exceptions"
	"forum/pkg/cust_encoders"
)

func (h *Handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		var ctx context.Context
		session, sessionErr := h.AuthHandler.AuthService.GetSession()
		fmt.Println(session)
		cookie, cookieErr := r.Cookie("session")
		if _, exclude := h.ExcludeSessionHandlersPath[r.URL.Path]; !exclude {
			if sessionErr != nil {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			if cookieErr != nil {
				newCookie := &http.Cookie{
					Name:     "session",
					Value:    session.Token,
					Path:     "/",
					Expires:  session.ExpireTime,
					HttpOnly: true,
					MaxAge:   7200,
				}
				http.SetCookie(w, newCookie)
				ctx := context.WithValue(r.Context(), "session", session)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if session.Token != cookie.Value {
				dataErr := exceptions.NewAuthenticationError()
				params := cust_encoders.EncodeParams(dataErr)
				http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			if sessionErr == nil {
				if r.URL.Path == "/signin" {
					fmt.Println("asdasds")
					if session.Token == cookie.Value {
						http.Redirect(w, r, "/post/", http.StatusSeeOther)
						return
					}
					err := h.AuthHandler.AuthService.DeleteSession()
					if err != nil {
						params := cust_encoders.EncodeParams(err)
						if err != nil {
							http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
							return
						}
					}
					ctx = r.Context()
				} else {
					ctx = context.WithValue(r.Context(), "session", session)
				}
			} else {
				ctx = r.Context()
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (h *Handler) ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			stringParams := r.URL.RawQuery
			if stringParams == "" {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
			} else {
				dataErr, err := cust_encoders.DecodeParams(stringParams)
				if err != nil {
					dataErr := exceptions.NewBadRequestError()
					params := cust_encoders.EncodeParams(dataErr)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				}
				var (
					customErr error
					ctx       context.Context
				)
				switch dataErr.(type) {
				case exceptions.AuthenticationError:
					customErr = dataErr
				case exceptions.ForbiddenError:
					customErr = dataErr
				case exceptions.ResourceNotFoundError:
					customErr = dataErr
				case exceptions.StatusMethodNotAllowed:
					customErr = dataErr
				case exceptions.StatusConflictError:
					customErr = dataErr
				case exceptions.ValidationError:
					customErr = dataErr
				case exceptions.InternalServerError:
					customErr = dataErr
				case exceptions.BadRequestError:
					customErr = dataErr
				default:
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}

				ctx = context.WithValue(r.Context(), "error", customErr)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			dataErr := exceptions.NewResourceNotFoundError()
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		}
	})
}
