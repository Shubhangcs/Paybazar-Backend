package interfaces

import "github.com/labstack/echo/v4"

type AuthInterface interface {
	RegisterAdmin(echo.Context) (string, error)
	RegisterMasterDistributor(echo.Context) (string, error)
	RegisterDistributor(echo.Context) (string, error)
	RegisterUser(echo.Context) (string, error)
}
