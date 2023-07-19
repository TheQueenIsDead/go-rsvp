package container

import (
	"database/sql"
	"github.com/labstack/echo/v4"
)

type Application struct {
	Server   *echo.Echo
	Database *sql.DB
}
