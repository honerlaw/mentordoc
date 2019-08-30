package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

const TokenRefresh = "refresh_token";

const TokenAccess = "access_token";

const tokenAccessExpireTime = time.Hour
const tokenRefreshExpireTime = 7 * 24 * time.Hour // 7 days
const issuer = "mentordoc"

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (service *TokenService) GenerateToken(resourceId string, tokenType string) (*string, error) {
	if tokenType != TokenRefresh && tokenType != TokenAccess {
		return nil, errors.New("invalid token type")
	}

	// change how long the token lives for based on the type of token we are isssuing
	timeUntilExpire := tokenAccessExpireTime
	if tokenType == TokenRefresh {
		timeUntilExpire = tokenRefreshExpireTime
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(timeUntilExpire).Unix(),
		IssuedAt:  time.Now().Unix(),
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

func (service *TokenService) ParseAndValidateToken(tokenValue string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		log.Print("failed to parse jwt", err)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims);
	if !ok || !token.Valid {
		log.Print("claims is not a standard claims struct or token is not valid")
		return nil, errors.New("invalid token")
	}

	// unsupported token type
	if claims.Audience != TokenAccess && claims.Audience != TokenRefresh {
		log.Print("invalid audience on JWT", claims.Audience)
		return nil, errors.New("invalid token")
	}

	// bad issuer
	if claims.Issuer != issuer {
		log.Print("invalid issuer on JWT", claims.Issuer)
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
