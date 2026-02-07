package utility

import (
	"context"
	"errors"
	"fmt"
	"realtime_quiz_system/internal/config"
	"realtime_quiz_system/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ConfigJWTHashFunction = jwt.SigningMethodHS256

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword compares a password with a hash
func VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenJWT(payload jwt.MapClaims) (token string, err error) {
	tokenJWT := jwt.NewWithClaims(ConfigJWTHashFunction, payload)
	secretKeyBytes := []byte(config.GetConfig().Auth.SecretKey)
	token, err = tokenJWT.SignedString(secretKeyBytes)
	if err != nil {
		g.Log().Error(context.Background(), "GenJWT error: ", err)
		return
	}
	return
}

// Validate JWT token
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	tokenParse, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().Auth.SecretKey), nil
	})

	var claims jwt.MapClaims

	if tokenParse != nil {
		var ok bool
		claims, ok = tokenParse.Claims.(jwt.MapClaims)
		if !ok {
			return nil, gerror.NewCode(consts.CodeInvalidToken)
		}
	}

	switch {
	case err == nil && tokenParse.Valid:
		return claims, nil
	case err != nil && (errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet)):
		return claims, jwt.ErrTokenExpired
	default:
		return nil, gerror.NewCode(consts.CodeInvalidToken)
	}
}
