package auth

import (
	"errors"
	"github.com/duxphp/duxgo/v2/exception"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

// Middleware 授权中间件
func Middleware(app string, renewals ...int64) echo.MiddlewareFunc {
	key := []byte(registry.Config["app"].GetString("app.safeKey"))
	// 续期时间
	var renewal int64 = 43200
	if len(renewals) > 0 {
		renewal = renewals[0]
	}
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: key,
		ParseTokenFunc: func(c echo.Context, token string) (interface{}, error) {
			data := jwt.MapClaims{}
			jwtToken, err := jwt.ParseWithClaims(token, &data, func(token *jwt.Token) (interface{}, error) {
				return key, nil
			})
			if err != nil {
				return nil, err
			}
			if sub, ok := data["sub"].(string); !ok || sub != app {
				return nil, errors.New("token type error")
			}
			return jwtToken.Valid, nil
		},
		SuccessHandler: func(c echo.Context) {
			// 获取token
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			c.Set("auth", cast.ToStringMap(claims))

			// 重新续期
			iat := claims["iat"].(int64) // 签发时间
			exp := claims["exp"].(int64) // 过期时间
			expire := exp - iat
			unix := time.Now().Unix()
			if iat+renewal <= unix {
				var refToken string
				refToken, _ = NewJWT().MakeToken(app, claims, expire)
				c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+refToken)
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return exception.New(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		},
	})
}
