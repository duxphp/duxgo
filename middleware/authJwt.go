package middleware

import (
	"errors"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/core/exception"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

// AuthJwt 授权中间件
func AuthJwt(auth string) echo.MiddlewareFunc {
	key := []byte(core.Config["app"].GetString("app.safeKey"))
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  key,
		TokenLookup: "header:" + echo.HeaderAuthorization + ",query:auth",
		Claims:      jwt.MapClaims{},
		ParseTokenFunc: func(token string, c echo.Context) (interface{}, error) {
			data := jwt.MapClaims{}
			jwtToken, err := jwt.ParseWithClaims(token, data, func(token *jwt.Token) (interface{}, error) {
				return key, nil
			})
			if err != nil {
				return nil, err
			}
			if userAuth, ok := data["auth"].(string); !ok || userAuth != auth {
				return nil, errors.New("token type error")
			}
			return jwtToken, nil
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			c.Set("authID", claims["id"])
			if claims["refresh"] == nil {
				return
			}
			if validity, ok := claims["refresh"].(float64); ok {
				tm := time.Unix(int64(validity), 0).Unix()
				if tm < time.Now().Unix() {
					var refToken string
					refToken, _ = NewJWT().MakeToken(auth, uint(claims["id"].(float64)), claims["extend"])
					c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+refToken)
				}
			}
		},
		ErrorHandler: func(err error) error {
			return exception.New(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		},
	})
}

// NewJWT 新建授权结构
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(core.Config["app"].GetString("app.safeKey")),
		Expired:    core.Config["app"].GetDuration("jwt.expired"), // 过期时间
		Renewal:    core.Config["app"].GetDuration("jwt.renewal"), // 续期时间
	}
}

// JWT 结构体
type JWT struct {
	SigningKey []byte
	Expired    time.Duration
	Renewal    time.Duration
}

// MakeToken 生成 token
func (j *JWT) MakeToken(auth string, userId uint, extend ...any) (tokenString string, err error) {

	var ext any
	if len(extend) > 0 {
		ext = extend[0]
	}

	claim := jwt.MapClaims{
		"auth":    auth,
		"id":      userId,
		"exp":     time.Now().Add((j.Expired + j.Renewal) * time.Minute).Unix(), // 过期时间
		"refresh": time.Now().Add(j.Expired * time.Minute).Unix(),               // 刷新时间
		"extend":  ext,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(j.SigningKey)
	return tokenString, err
}

// ParsingToken 解析 token
func (j *JWT) ParsingToken(auth string, token string) (claims jwt.MapClaims, err error) {
	data := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if data["auth"] != auth {
		return nil, errors.New("token type error")
	}
	return data, nil
}
