package server

import (
	"context"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"net/http"
	"strings"
)

const AuthenticatedUserContextKey = "authenticated_user"

type AuthenticationMiddleware struct {
	authenticationService *AuthenticationService
	userService           *UserService
}

func NewAuthenticationMiddleware(authenticationService *AuthenticationService, userService *UserService) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authenticationService: authenticationService,
		userService:           userService,
	}
}

// @todo create one of these for the refresh token
func (middleware *AuthenticationMiddleware) HasAccessToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			header := req.Header.Get("Authorization")
			if header == "" {
				util.WriteHttpError(w, model.NewUnauthorizedError("invalid token"))
				return
			}

			pieces := strings.Split(header, "Bearer ")
			token := pieces[0]

			claims, err := middleware.authenticationService.ParseAndValidateToken(token)
			if err != nil {
				util.WriteHttpError(w, model.NewUnauthorizedError("invalid token"))
				return
			}

			// make sure they are using an access token
			if claims.Audience != TokenAccess {
				util.WriteHttpError(w, model.NewUnauthorizedError("invalid token"))
				return
			}

			// lookup the user
			user := middleware.userService.FindById(claims.Subject)
			if user == nil {
				util.WriteHttpError(w, model.NewUnauthorizedError("invalid token"))
				return
			}

			// store the user on the request context
			ctx := context.WithValue(req.Context(), AuthenticatedUserContextKey, user)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func (middleware *AuthenticationMiddleware) GetUserFromRequest(req *http.Request) *model.User {
	return req.Context().Value(AuthenticatedUserContextKey).(*model.User)
}
