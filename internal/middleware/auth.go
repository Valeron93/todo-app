package middleware

import (
	"net/http"

	"github.com/Valeron93/todo-app/internal/model"
)

type AuthMiddleware struct {
	sessions model.SessionManager
}

func NewAuthMiddleware(sessionManager model.SessionManager) AuthMiddleware {
	return AuthMiddleware{
		sessions: sessionManager,
	}
}

func (a *AuthMiddleware) Unauthorized401(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, ok := model.SessionFromCtx(r.Context())

		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *AuthMiddleware) UnauthorizedRedirect(url string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			_, ok := model.SessionFromCtx(r.Context())

			if !ok {
				http.Redirect(w, r, url, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *AuthMiddleware) InjectSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if sessionCookie, err := r.Cookie("session_token"); err == nil {
			if session, err := a.sessions.GetSession(sessionCookie.Value); err == nil {
				ctx = model.CtxWithSession(ctx, session)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
