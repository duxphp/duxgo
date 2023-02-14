package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/duxphp/duxgo/v2/function"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"net/url"
	"strings"
	"time"
)

// 误差秒
const diffTime float64 = 5

// Middleware 接口签名
func Middleware(secretCallback func(id string) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			date := c.Request().Header.Get("Content-Date")
			timeNow := time.Now()
			t := time.Unix(cast.ToInt64(date), 0)
			if timeNow.Sub(t).Seconds() > diffTime {
				return echo.ErrRequestTimeout
			}
			// 签名验证
			sign := c.Request().Header.Get("Content-MD5")
			id := c.Request().Header.Get("AccessKey")

			secretKey := secretCallback(id)
			if secretKey == "" {
				return echo.ErrUnauthorized
			}
			body := function.CtxBody(c)
			md5 := strings.ToLower(function.Md5(string(body)))
			query, _ := url.QueryUnescape(c.Request().URL.RawQuery)
			signData := []string{
				c.Request().URL.Path,
				query,
				md5,
				date,
			}
			h := sha256.New
			mac := hmac.New(h, []byte(secretKey))
			mac.Write([]byte(strings.Join(signData, "\n")))
			digest := mac.Sum(nil)
			hexDigest := hex.EncodeToString(digest)
			if sign != hexDigest {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
