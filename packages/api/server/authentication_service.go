package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/honerlaw/mentordoc/server/util"
	"github.com/pkg/errors"
	"log"
	"os"
)

const TokenRefresh = "refresh_token";

const TokenAccess = "access_token";

const tokenAccessExpireTime int64 = 60 * 60 * 1000           // 1 hour
const tokenRefreshExpireTime int64 = 7 * 24 * 60 * 60 * 1000 // 1 week
const issuer = "mentordoc"

type AuthenticationService struct{}

func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{}
}

func (service *AuthenticationService) GenerateToken(resourceId string, tokenType string) (*string, error) {
	if tokenType != TokenRefresh && tokenType != TokenAccess {
		return nil, errors.New("invalid token type")
	}

	// change how long the token lives for based on the type of token we are issuing
	timeUntilExpire := tokenAccessExpireTime
	if tokenType == TokenRefresh {
		timeUntilExpire = tokenRefreshExpireTime
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: util.NowUnix() + timeUntilExpire,
		IssuedAt:  util.NowUnix(),
		Issuer:    issuer,
		Subject:   resourceId,
		Audience:  tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenValue, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to sign token")
	}

	return &tokenValue, nil
}

func (service *AuthenticationService) ParseAndValidateToken(tokenValue string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		log.Print(err)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims);
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// unsupported token type
	if claims.Audience != TokenAccess && claims.Audience != TokenRefresh {
		return nil, errors.New("invalid token")
	}

	// bad issuer
	if claims.Issuer != issuer {
		return nil, errors.New("invalid token")
	}

	// expired
	if util.NowUnix()-claims.ExpiresAt > 0 {
		return nil, errors.New("expired token")
	}

	return claims, nil
}
