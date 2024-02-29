package handler

import (
	"context"
	"forum/internal/exceptions"
	"forum/pkg/cust_encoders"
	"net/http"
)

func (h *Handler) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session, sessionErr := h.AuthHandler.AuthService.GetSession()
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
				expiredCookie := &http.Cookie{
					Name:     "session",
					Value:    "",
					Path:     "/",
					HttpOnly: true,
					MaxAge:   -1,
				}
				http.SetCookie(w, expiredCookie)

				http.Redirect(w, r, "/post/", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			if sessionErr == nil {
				if r.URL.Path == "/signin" {
					if cookieErr != nil {
						h.AuthHandler.AuthService.DeleteSession()
					} else if session.Token == cookie.Value {
						http.Redirect(w, r, "/post/", http.StatusSeeOther)
						return
					} else {
						err := h.AuthHandler.AuthService.DeleteSession()
						if err != nil {
							params := cust_encoders.EncodeParams(err)
							if err != nil {
								http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
								return
							}
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
					dataErr := exceptions.NewBadRequestError("Invalid error parameters are given")
					params := cust_encoders.EncodeParams(dataErr)
					http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
				}
				var (
					customErr  error
					ctx        context.Context
					httpStatus int
				)
				switch dataErr.(type) {
				case exceptions.AuthenticationError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.AuthenticationError).StatusCode
				case exceptions.ForbiddenError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.ForbiddenError).StatusCode
				case exceptions.ResourceNotFoundError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.ResourceNotFoundError).StatusCode
				case exceptions.StatusMethodNotAllowed:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.StatusMethodNotAllowed).StatusCode
				case exceptions.StatusConflictError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.StatusConflictError).StatusCode
				case exceptions.ValidationError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.ValidationError).StatusCode
				case exceptions.InternalServerError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.InternalServerError).StatusCode
				case exceptions.BadRequestError:
					customErr = dataErr
					httpStatus = dataErr.(exceptions.BadRequestError).StatusCode
				default:
					http.Redirect(w, r, "/signin", http.StatusSeeOther)
					return
				}

				ctx = context.WithValue(r.Context(), "error", customErr)
				ctx = context.WithValue(ctx, "httpStatus", httpStatus)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			dataErr := exceptions.NewResourceNotFoundError("Page is not found")
			params := cust_encoders.EncodeParams(dataErr)
			http.Redirect(w, r, "/?"+params, http.StatusSeeOther)
		}
	})
}
