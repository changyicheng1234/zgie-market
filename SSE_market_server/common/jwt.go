// 返回token的文件，不用管，要看也行。模板写法来着，至于为什么，我也不是很懂
package common

import (
	"loginTest/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// 采用jwt方式生成token
var jwtKey []byte

type Claims struct {
	UserID int
	jwt.StandardClaims
}

type Claims_admin struct {
	AdminID int
	jwt.StandardClaims
}

// RefreshClaims 用于刷新 token 的 JWT 声明
type RefreshClaims struct {
	UserID int
	jwt.StandardClaims
}

// OAuth2Claims OAuth2访问令牌的Claims
type OAuth2Claims struct {
	UserID int
	AppID  string
	Scope  string
	jwt.StandardClaims
}

func InitJWTkey() {
	jwtKey = []byte(viper.GetString("crypto.jwtKey"))
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "loginTest",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ReleaseRefreshToken 生成用于刷新 access token 的 refresh token
// 有效期：30 天
func ReleaseRefreshToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &RefreshClaims{
		UserID: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "loginTest",
			Subject:   "refresh token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ReleaseToken_admin(admin model.Admin) (string, error) {
	expirationTime := time.Now().Add(1 * 24 * time.Hour)
	claims := &Claims_admin{
		AdminID: admin.AdminID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "loginTest",
			Subject:   "admin token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

func ParseToken_admin(tokenString string) (*jwt.Token, *Claims_admin, error) {
	claims_admin := &Claims_admin{}
	token, err := jwt.ParseWithClaims(tokenString, claims_admin, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims_admin, err
}

// ParseRefreshToken 解析 refresh token
func ParseRefreshToken(tokenString string) (*jwt.Token, *RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// 生成OAuth2访问令牌
func ReleaseOAuth2Token(userID int, appID, scope string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // 有效期5分钟
	claims := &OAuth2Claims{
		UserID: userID,
		AppID:  appID,
		Scope:  scope,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "loginTest",
			Subject:   "oauth2 token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析OAuth2访问令牌
func ParseOAuth2Token(tokenString string) (*jwt.Token, *OAuth2Claims, error) {
	claims := &OAuth2Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
