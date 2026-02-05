package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/skeleton-golang/internal/infrastructure/container"
)

func SetupRouter(e *echo.Echo, cnt *container.Container) {
	h := SetupHandler(cnt).Validate()

	v1 := e.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.GET("", h.userHandler.GetAllUser)
		}
	}

}
