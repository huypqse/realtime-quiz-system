package service

import (
	"context"
	"encoding/json"
	"fmt"
	"realtime_quiz_system/internal/config"
	"realtime_quiz_system/internal/consts"
	"realtime_quiz_system/utility"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
)

// TokenService handles token generation and verification
type TokenService interface {
	NewAccessToken(userId, username string) *AccessToken
	NewRefreshToken(ctx context.Context, uuid, uid string) *RefreshToken
	GenerateAuthToken(userId, username string) (string, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiresAt int64  `json:"expires_at"`
}

type tokenService struct {
	config *config.Config
}

// NewTokenService creates a new token service instance
func NewTokenService(cfg *config.Config) TokenService {
	return &tokenService{
		config: cfg,
	}
}

// GenerateAuthToken generates JWT token with 7-day expiration
func (ts *tokenService) GenerateAuthToken(userId, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // 7 days

	claims := jwt.MapClaims{
		"user_id":    userId,
		"username":   username,
		"issued_at":  now.Unix(),
		"expires_at": expiresAt.Unix(),
		"exp":        expiresAt.Unix(), // Standard JWT expiration claim
		"iat":        now.Unix(),       // Standard JWT issued at claim
	}

	return utility.GenJWT(claims)
}

// ValidateToken validates JWT token and returns claims
func (ts *tokenService) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	jwtMap, err := utility.ParseJWT(token)
	if err != nil {
		return nil, err
	}

	userId, ok := jwtMap["user_id"].(string)
	if !ok {
		return nil, gerror.NewCode(consts.CodeInvalidToken)
	}

	username, ok := jwtMap["username"].(string)
	if !ok {
		return nil, gerror.NewCode(consts.CodeInvalidToken)
	}

	issuedAt, ok := jwtMap["issued_at"].(float64)
	if !ok {
		return nil, gerror.NewCode(consts.CodeInvalidToken)
	}

	expiresAt, ok := jwtMap["expires_at"].(float64)
	if !ok {
		return nil, gerror.NewCode(consts.CodeInvalidToken)
	}

	return &TokenClaims{
		UserId:    userId,
		Username:  username,
		IssuedAt:  int64(issuedAt),
		ExpiresAt: int64(expiresAt),
	}, nil
}

func (ts *tokenService) NewAccessToken(userId, username string) *AccessToken {
	return &AccessToken{
		Iss:    userId,
		Sub:    username,
		config: ts.config,
	}
}

func (ts *tokenService) NewRefreshToken(ctx context.Context, uuid, uid string) *RefreshToken {
	return &RefreshToken{
		Ctx:    ctx,
		Uuid:   uuid,
		Uid:    uid,
		config: ts.config,
	}
}

type Token interface {
	Gen() (token string, err error)
	Verify(token string) (err error)
}

type AccessToken struct {
	Iss    string `json:"iss"`
	Sub    string `json:"sub"`
	Exp    int64  `json:"exp"`
	config *config.Config
}

func (ac *AccessToken) Gen() (token string, err error) {
	expTime := gtime.Now().Add(time.Duration(ac.config.Auth.AccessTokenExpireMinute) * time.Minute)
	ac.Exp = expTime.Unix()

	dataByte, err := json.Marshal(ac)
	if err != nil {
		return "", err
	}
	var mapClaims jwt.MapClaims
	if err := json.Unmarshal(dataByte, &mapClaims); err != nil {
		return "", err
	}
	return utility.GenJWT(mapClaims)
}

func (ac *AccessToken) Verify(ctx context.Context, token string) (err error) {
	jwtMap, err := utility.ParseJWT(token)
	if err != nil {
		return err
	}
	dataByte, _ := json.Marshal(jwtMap)
	err = json.Unmarshal(dataByte, &ac)
	if err != nil {
		return err
	}
	return
}

type RefreshToken struct {
	Ctx    context.Context
	Uuid   string
	Uid    string
	Exp    string
	Nbf    string
	config *config.Config
}

func (rf *RefreshToken) CheckCtx() (err error) {
	if rf.Ctx == nil {
		return gerror.NewCode(consts.CodeInvalidToken)
	}
	return
}

func (rf *RefreshToken) Gen() (token string, err error) {
	expTime := gtime.Now().Add(time.Duration(rf.config.Auth.RefreshTokenExpireMinute) * time.Minute)
	rf.Exp = expTime.String()
	rf.Nbf = fmt.Sprintf("%v", expTime.Unix())
	token = fmt.Sprintf("%v", rf.Uuid)
	return
}

func (rf *RefreshToken) Verify(token string) (err error) {
	if gtime.Now().Unix() >= utility.String2Int64(rf.Nbf) {
		return gerror.NewCode(consts.CodeTokenExpired)
	}
	return
}
