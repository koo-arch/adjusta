package auth

import (
	"context"
	"strings"
	"time"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

var (
	accessTokenTTL = 15 * time.Minute // アクセストークンの有効期限
	refreshTokenTTL = 24 * time.Hour  // リフレッシュトークンの有効期限
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type JWTManager struct {
	client *ent.Client
	keyManager *KeyManager
}

func NewJWTManager(client *ent.Client, km *KeyManager) *JWTManager {
	return &JWTManager{
		client: client,
		keyManager: km,
	}
}

// GenerateAccessToken generates a new JWT access token
func (jm *JWTManager) GenerateAccessToken(ctx context.Context, client *ent.Client, email string) (string, time.Time, error) {
	expirationTime := time.Now().Add(accessTokenTTL)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// アクセスキーの取得

	accessKey, err := jm.keyManager.GetJWTKey(ctx, "access"); if err != nil {
		log.Println("failed to get access key")
		return "", time.Time{}, err
	}

	signedToken, err := token.SignedString(accessKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expirationTime, nil
}

// GenerateRefreshToken generates a new JWT refresh token
func (jm *JWTManager) GenerateRefreshToken(ctx context.Context, client *ent.Client, email string) (string, time.Time, error) {
	expirationTime := time.Now().Add(refreshTokenTTL)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// リフレッシュキーの取得
	refreshKey, err := jm.keyManager.GetJWTKey(ctx, "refresh"); if err != nil {
		log.Println("failed to get refresh key")
		return "", time.Time{}, err
	}

	signedToken, err := token.SignedString(refreshKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expirationTime, nil
}

// GenerateTokens generates a new JWT access and refresh token pair
func (jm *JWTManager) GenerateTokens(ctx context.Context, client *ent.Client, email string) (*models.JWTToken, error) {
	accessToken, accessExpiration, err := jm.GenerateAccessToken(ctx, client, email)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExpiration, err := jm.GenerateRefreshToken(ctx, client, email)
	if err != nil {
		return nil, err
	}

	jwtToken := &models.JWTToken{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		AccessExpiration: accessExpiration,
		RefreshExpiration: refreshExpiration,
	}

	return jwtToken, nil
}

// VerifyToken verifies a JWT token and returns the email
func (jm *JWTManager) VerifyToken(ctx context.Context, client *ent.Client, tokenString, tokenType string) (string, error) {
	// JWTキーの取得
	keys, err := jm.keyManager.GetJWTKeys(ctx, tokenType); if err != nil {
		log.Println("failed to get access key")
		return "", err
	}

	for _, key := range keys {
		email, err := jm.verifyTokenWithKey(tokenString, []byte(key.Key))
		if err == nil {
			return email, nil
		} else if strings.Contains(err.Error(), "token is expired") {
			return "", err
		}
	}

	return "", jwt.ErrSignatureInvalid
}

func (jm *JWTManager) verifyTokenWithKey(tokenString string, key []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", jwt.ErrSignatureInvalid
		}
		return key, nil
	})
	if err != nil {
		println("failed to parse token", err.Error())
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		println("invalid token")
		return "", jwt.ErrSignatureInvalid
	}

	return claims.Email, nil
}