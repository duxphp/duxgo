package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// NewJWT 新建授权结构
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(config.Get("app").GetString("app.safeKey")),
	}
}

// JWT 结构体
type JWT struct {
	SigningKey []byte
}

// MakeToken 生成 token
func (j *JWT) MakeToken(app string, params jwt.MapClaims, expires ...int64) (tokenString string, err error) {
	var expire int64 = 86400
	if len(expires) > 0 {
		expire = expires[0]
	}
	claim := jwt.MapClaims{
		"sub": app,
		"exp": time.Now().Add(time.Duration(expire) * time.Minute).Unix(), // 过期时间
		"iat": time.Now().Unix(),                                          // 签发时间
	}
	for key, value := range params {
		claim[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(j.SigningKey)
	return tokenString, err
}

// ParsingToken 解析 token
func (j *JWT) ParsingToken(app string, token string) (claims jwt.MapClaims, err error) {
	data := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if data["sub"] != app {
		return nil, errors.New("token type error")
	}
	return data, nil
}
