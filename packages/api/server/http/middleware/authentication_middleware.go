package middleware

import (
	"context"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/user"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
	"net/http"
	"strings"
)

const AuthenticatedUserContextKey = "authenticated_user"

type AuthenticationMiddleware struct {
	tokenService *util.TokenService
	userService  *user.UserService
}

func NewAuthenticationMiddleware(tokenService *util.TokenService, userService *user.UserService) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		tokenService: tokenService,
		userService:  userService,
	}
}

// @todo create one of these for the refresh token
func (middleware *AuthenticationMiddleware) HasAccessToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			header := req.Header.Get("Authorization")
			if header == "" {
				util.WriteHttpError(w, shared.NewUnauthorizedError("invalid token"))
				return
			}

			pieces := strings.Split(header, "Bearer ")
			token := pieces[1]

			claims, err := middleware.tokenService.ParseAndValidateToken(token)
			if err != nil {
				util.WriteHttpError(w, shared.NewUnauthorizedError("invalid token"))
				return
			}

			// make sure they are using an access token
			if claims.Audience != util.TokenAccess {
				log.Print("attempted to use refresh token instead of access token")
				util.WriteHttpError(w, shared.NewUnauthorizedError("invalid token"))
				return
			}

			// lookup the user
			user := middleware.userService.FindById(claims.Subject)
			if user == nil {
				log.Print("failed to find the user for subject", claims.Subject)
				util.WriteHttpError(w, shared.NewUnauthorizedError("invalid token"))
				return
			}

			// store the user on the request context
			ctx := context.WithValue(req.Context(), AuthenticatedUserContextKey, user)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func (middleware *AuthenticationMiddleware) GetUserFromRequest(req *http.Request) *shared.User {
	return req.Context().Value(AuthenticatedUserContextKey).(*shared.User)
}
