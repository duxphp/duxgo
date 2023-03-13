package auth

import (
	"github.com/demdxx/gocast/v2"
	"github.com/duxphp/duxgo/v2/config"
	"github.com/gofiber/fiber/v2"
	fiberJwt "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Middleware 授权中间件
func Middleware(app string, renewals ...int64) fiber.Handler {
	key := []byte(config.Get("app").GetString("app.safeKey"))
	// 续期时间
	var renewal int64 = 43200
	if len(renewals) > 0 {
		renewal = renewals[0]
	}

	return fiberJwt.New(fiberJwt.Config{
		SigningKey: key,
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			c.Locals("auth", gocast.Map[string, any](claims))
			// 判断应用
			sub, ok := claims["sub"].(string)
			if !ok || sub != app {
				return c.Status(fiber.StatusUnauthorized).SendString("token type error jwt")
			}
			// 验证刷新
			iat := claims["iat"].(int64) // 签发时间
			exp := claims["exp"].(int64) // 过期时间
			unix := time.Now().Unix()
			if iat+renewal <= unix {
				expire := exp - iat
				claims["exp"] = unix + expire
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString(token)
				c.Set(fiber.HeaderAuthorization, "Bearer "+tokenString)
			}
			return c.Next()
		},
	})
}
