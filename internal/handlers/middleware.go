package handler

import (
	"context"
	"fmt"
	"net/http"
)

func (h *Handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			// log.Fatal(err)
		}
		fmt.Println(cookie.Name)
		fmt.Println("Value", cookie.Value)

		session, err := h.AuthHandler.AuthService.GetSession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		fmt.Println(session.UserID)

		ctx := context.WithValue(r.Context(), "user_id", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
