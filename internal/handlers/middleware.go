package handler

import (
	"context"
	"fmt"
	"forum/internal/exceptions"
	"forum/pkg/cust_encoders"
	"net/http"
)

func (h *Handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		session, err := h.AuthHandler.AuthService.GetSession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		fmt.Println("Sessoin is contex is set")

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			errorStringParam := r.URL.Query().Get("params")
			if errorStringParam == "" {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
			} else {
				err := cust_encoders.DecodeParams(errorStringParam)
				var (
					customErr error
					ctx       context.Context
				)
				switch err.(type) {
				case exceptions.AuthenticationError:
					customErr = err
				case exceptions.ForbiddenError:
					customErr = err
				case exceptions.ResourceNotFoundError:
					customErr = err
				case exceptions.StatusMethodNotAllowed:
					customErr = err
				case exceptions.StatusConflictError:
					customErr = err
				case exceptions.ValidationError:
					customErr = err
				case exceptions.InternalServerError:
					customErr = err
				case exceptions.BadRequestError:
					customErr = err
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
			http.Redirect(w, r, "/?params="+params, http.StatusSeeOther)
		}
	})
}
