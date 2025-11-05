package interfaces

import "github.com/labstack/echo/v4"

type AuthInterface interface {
	RegisterAdmin(echo.Context) (string,error)
}