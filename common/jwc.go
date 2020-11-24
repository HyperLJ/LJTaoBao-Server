package common

import (
	"github.com/dgrijalva/jwt-go"
	"itStudioTB/model"
	"time"
)

// 生成token的密钥
var jwtKey = []byte("a_taobao_secret")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 发布token
func IssueToken(user model.User) (string, error) {
	// token 过期时间 7天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserId: user.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // token 过期时间
			IssuedAt:  time.Now().Unix(),     // token 发放时间
			Issuer:    "TaoBaoServer",        // token 发放人
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey) // 根据密钥生成token

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	// token组成:编码方式 + claims + jwt密钥

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
