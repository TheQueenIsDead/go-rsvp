package container

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Application struct {
	Server   *echo.Echo
	Database *gorm.DB
}
