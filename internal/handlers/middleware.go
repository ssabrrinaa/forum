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
			//errorStringParam := r.URL.Query().Get("params")
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
