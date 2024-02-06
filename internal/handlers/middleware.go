package handler

import (
	"context"
	"fmt"
	"forum/internal/exceptions"
	"net/http"
	"strconv"
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
		errorStringCode := r.URL.Query().Get("error")
		var ctx context.Context
		var customErr error
		if errorStringCode != "" && r.URL.Path == "/" {
			code, _ := strconv.Atoi(errorStringCode)
			customErr = exceptions.NewBadRequestError()
			switch code {
			case 401:
				customErr = exceptions.NewAuthenticationError()
			case 403:
				customErr = exceptions.NewForbiddenError()
			case 404:
				customErr = exceptions.NewResourceNotFoundError()
			case 405:
				customErr = exceptions.NewStatusMethodNotAllowed()
			case 409:
				customErr = exceptions.NewStatusConflicError()
			case 422:
				customErr = exceptions.NewValidationError()
			case 500:
				customErr = exceptions.NewInternalServerError()
			}
			ctx = context.WithValue(r.Context(), "error", customErr)
		} else {
			http.Redirect(w, r, "/?error=404", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
