package permission

import (
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"sync"
)

// Middleware 权限中间件
func Middleware(app string, model any) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		var doOnce sync.Once
		var registeredRoutes []*echo.Route
		var permissions []string
		return func(c echo.Context) error {
			auth, ok := c.Get("auth").(map[string]any)
			if !ok {
				return echo.ErrUnauthorized
			}
			doOnce.Do(func() {
				registeredRoutes = c.Echo().Routes()
				permissions = Get(app).GetData()
			})
			routeName := ""
			for _, r := range registeredRoutes {
				if r.Method == c.Request().Method && r.Path == c.Path() {
					routeName = r.Name
					break
				}
			}
			if routeName == "" || lo.IndexOf[string](permissions, routeName) == -1 {
				return next(c)
			}
			info := map[string]any{}
			err := registry.Db.Model(model).Where("id = ?", auth["id"]).Find(info).Error
			if err != nil {
				return err
			}
			permission := cast.ToStringSlice(info["permission"])
			if len(permission) > 0 && lo.IndexOf[string](permission, routeName) == -1 {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}
